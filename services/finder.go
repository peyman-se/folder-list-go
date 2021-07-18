package services

import (
	"os"
	"sync"
	"sort"
	"log"
	"example.com/finder/models"
)

var dirPath string
func GetFiles(dirPath string) ([]os.FileInfo, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return []os.FileInfo{}, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return []os.FileInfo{}, err
	}

	return files, err
}

func GetEntitiesOrderedBySizeFromPath(path string) ([]models.FinderEntity, int64, error) {
	dirPath = path
	files, err := GetFiles(dirPath)
	if err != nil {
		return []models.FinderEntity{}, int64(0), err
	}

	in := entityChan(files)

	//build as many routines as is appropriate
	c1, e1 := getSizeChan(in)
	c2, e2 := getSizeChan(in)
	c3, e3 := getSizeChan(in)

	// if any of routines returned an error, signal other routines to stop maybe? Depends on business logic
	if len(e1) != 0 || len(e2) !=0 || len(e3) != 0 {
		log.Println("calculation is not complete due to an error.")
	}
	var totalSize int64
	var entities []models.FinderEntity

	for entity := range merge(c1, c2, c3) {
		entities = append(entities, entity)
		totalSize += entity.Size
	}
	
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Size < entities[j].Size
	})

	return entities, totalSize, nil
}

func entityChan(files []os.FileInfo) <-chan models.FinderEntity {
	out := make(chan models.FinderEntity)
	go func() {
		for _, file := range files {
			size := file.Size()
			entity := models.FinderEntity{
				file,
				file.Name(),
				size,
				models.GetHumanReadableSize(size),
				dirPath,
				file.ModTime().Local().Format("Mon Jan 2 15:04:05 MST 2006"),
			}

			out <- entity
		}
		close(out)
	}()

	return out
}

func getSizeChan(inChan <-chan models.FinderEntity) (<-chan models.FinderEntity, <-chan error) {
	outChan := make(chan models.FinderEntity)
	errChan := make(chan error)
	go func() {
		for entity := range inChan {
			err := entity.SetActualSize()
			if err != nil {
				errChan <- err
			}
			
			outChan <- entity
		}

		close(outChan)
		close(errChan)
	}()

	return outChan, errChan
}

func merge(cs ...<-chan models.FinderEntity) <-chan models.FinderEntity {
	var wg sync.WaitGroup
	out := make(chan models.FinderEntity)

	output := func(c <-chan models.FinderEntity) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}