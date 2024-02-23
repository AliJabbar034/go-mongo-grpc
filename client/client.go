package main

import (
	"context"
	"fmt"

	blogpb "githu.com/alijabbar034/mongo_grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Client...")

	cc, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	blog := &blogpb.Blog{

		AuthorId: "alijabbar034",
		Title:    "Hello World",
		Content:  "Hello World",
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	created, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	fmt.Println(created)

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		Id: created.GetBlog().GetId(),
	})
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done")
}
