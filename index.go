// main.go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/handlers"
    "github.com/joho/godotenv"
)

// Response structure
type UserResponse struct {
    ID          string `json:"id"`
    Username    string `json:"username"`
    GlobalName  string `json:"global_name"`
    Assets      struct {
        AvatarURL string `json:"avatar_url"`
        BannerURL string `json:"banner_url"`
    } `json:"assets"`
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file")
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
        log.Println("Hello World")
    })

    http.HandleFunc("/discord/", func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Path[len("/discord/"):]

        token := os.Getenv("TOKEN")
        log.Println("Using Token:", token)

        req, err := http.NewRequest("GET", "https://canary.discord.com/api/v10/users/"+id, nil)
        if err != nil {
            log.Println("Error creating request:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        req.Header.Set("Authorization", "Bot "+token)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            log.Println("Error fetching user data:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        var jsonResponse map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
            log.Println("Error decoding JSON:", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        if message, ok := jsonResponse["message"]; ok {
            json.NewEncoder(w).Encode(message)
            return
        }

        userResponse := UserResponse{
            ID:         jsonResponse["id"].(string),
            Username:   jsonResponse["username"].(string),
            GlobalName: jsonResponse["global_name"].(string),
        }

        if avatar, ok := jsonResponse["avatar"]; ok {
            userResponse.Assets.AvatarURL = "https://cdn.discordapp.com/avatars/" + id + "/" + avatar.(string) + ".png"
        }
        if banner, ok := jsonResponse["banner"]; ok {
            userResponse.Assets.BannerURL = "https://cdn.discordapp.com/banners/" + id + "/" + banner.(string) + ".png"
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(userResponse)
    })

    log.Println("✓ [SERVER] Running on http://127.0.0.1:3000")
    log.Fatal(http.ListenAndServe("127.0.0.1:3000", handlers.CORS()(http.DefaultServeMux)))
}