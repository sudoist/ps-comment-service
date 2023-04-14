package controllers

import (
    "context"
    "gin-mongo-api/configs"
    "gin-mongo-api/models"
    "gin-mongo-api/responses"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")
var validate = validator.New()

func CreateComment() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var comment models.Comment
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&comment); err != nil {
            c.JSON(http.StatusBadRequest, responses.
            CommentResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&comment); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.CommentResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        newComment := models.Comment{
            ContentId:     comment.ContentId,
            UserId: comment.UserId,
            Comment: comment.Comment,
        }

        result, err := commentCollection.InsertOne(ctx, newComment)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusCreated, responses.CommentResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}

func GetAComment() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        commentId := c.Param("commentId")
        var comment models.Comment
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(commentId)

        err := commentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&comment)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusOK, responses.CommentResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": comment}})
    }
}

func EditAComment() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        commentId := c.Param("commentId")
        var comment models.Comment
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(commentId)

        //validate the request body
        if err := c.BindJSON(&comment); err != nil {
            c.JSON(http.StatusBadRequest, responses.CommentResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&comment); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.CommentResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        update := bson.M{"comment": comment.Comment}
        result, err := commentCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //get updated comment details
        var updatedComment models.Comment
        if result.MatchedCount == 1 {
            err := commentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedComment)
            if err != nil {
                c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
                return
            }
        }

        c.JSON(http.StatusOK, responses.CommentResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedComment}})
    }
}

func DeleteAComment() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        commentId := c.Param("commentId")
        defer cancel()

        objId, _ := primitive.ObjectIDFromHex(commentId)

        result, err := commentCollection.DeleteOne(ctx, bson.M{"_id": objId})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        if result.DeletedCount < 1 {
            c.JSON(http.StatusNotFound,
                responses.CommentResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Comment with specified ID not found!"}},
            )
            return
        }

        c.JSON(http.StatusOK,
            responses.CommentResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Comment successfully deleted!"}},
        )
    }
}

func GetAllComments() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var comments []models.Comment
        defer cancel()

        results, err := commentCollection.Find(ctx, bson.M{})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singleComment models.Comment
            if err = results.Decode(&singleComment); err != nil {
                c.JSON(http.StatusInternalServerError, responses.CommentResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }

            comments = append(comments, singleComment)
        }

        c.JSON(http.StatusOK,
            responses.CommentResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": comments}},
        )
    }
}
