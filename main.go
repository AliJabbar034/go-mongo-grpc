package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	blogpb "githu.com/alijabbar034/mongo_grpc/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var coll *mongo.Collection

type Server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	Id       primitive.ObjectID `bson:"_id"`
	AuthorId string             `bson:"author"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("BLOG server started..")

	uri := "mongodb+srv://mongoUser:mongo034@cluster0.9vklrcp.mongodb.net/sample_mflix?retryWrites=true&w=majority&appName=Cluster0"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll = client.Database("sample_mflix").Collection("movies")

	lis, err := net.Listen("tcp", "0.0.0.0:5051")
	if err != nil {

		fmt.Println(err)
		return
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(grpcServer, &Server{})

	go func() {
		fmt.Println("starting server...")
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Println(err)
			return
		}

	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	client.Disconnect(context.Background())
	fmt.Println("shutting down")
	grpcServer.GracefulStop()
	fmt.Println("shut down")
	lis.Close()

}

func (s *Server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := req.GetBlog()
	fmt.Println(blog)
	data := blogItem{
		AuthorId: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}
	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "insert failed")
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println(result.InsertedID)
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       id,
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil

}

func (s *Server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	id := req.GetId()
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	filter := bson.M{"_id": _id}
	result := &blogItem{}

	errr := coll.FindOne(ctx, filter).Decode(&result)

	if errr != nil {
		return nil, status.Errorf(codes.Internal, "read failed")
	}
	OId := result.Id.Hex()
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       OId,
			AuthorId: result.AuthorId,
			Title:    result.Title,
			Content:  result.Content,
		},
	}, nil
}
