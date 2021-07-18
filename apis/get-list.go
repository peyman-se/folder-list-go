package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"example.com/finder/models"
	"example.com/finder/services"
)

func GetList(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	// get user from token, for testing puposes we assume the user is always `hello`
	user := models.User{UserName: "hello"}

	// if path has a ~ then we should build the absolute path based on it
	// we should check permissions of the user as well, if user is not root
	// we should limit the access based on permissions
	if strings.Contains(path, "~") {
		path = strings.Replace(path, "~", "/Users/"+user.UserName, 1)
	}

	entities, totalSize, err := services.GetEntitiesOrderedBySizeFromPath(path)
	if err != nil {
		log.Printf("handle error properly: %v", err)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(
		map[string]interface{}{"entities": entities, "totalSize": models.GetHumanReadableSize(totalSize)},
	)
}
