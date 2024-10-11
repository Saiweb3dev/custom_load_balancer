package controllers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "simple_load_balancer/internal/models"
)

type UserController struct {
    collection *mongo.Collection
}

func NewUserController(db *mongo.Database) *UserController {
    return &UserController{
        collection: db.Collection("users"),
    }
}

func (uc *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := uc.collection.InsertOne(ctx, user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"id": result.InsertedID})
}

func (uc *UserController) GetLastUser(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := uc.collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"_id": -1})).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            http.Error(w, "No users found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(user)
}