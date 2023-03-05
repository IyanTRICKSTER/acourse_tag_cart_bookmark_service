package grpc_client

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"acourse_tag_cart_bookmark_service/pkg/models"
	ps "acourse_tag_cart_bookmark_service/pkg/models/proto_schema"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

type GRPCServiceClient struct {
	HOST   string
	PORT   string
	Client ps.CoursesServiceClient
}

func (c *GRPCServiceClient) List(ctx context.Context, coursesID []string) []models.Course {

	//cID := ps.CoursesID{CoursesID: []string{"6300988647b1637e7974b3d9", "6300988647b1637e7974b3d6"}}
	cID := ps.CoursesID{CoursesID: coursesID}

	courses, err := c.Client.List(ctx, &cID)
	if err != nil {
		log.Println("gRPC Client: CourseService: List Error >>", err)
		log.Fatal(err)
	}

	coursesResult := make([]models.Course, 0)

	for _, c := range courses.List {
		coursesResult = append(coursesResult, models.Course{ID: models.GenerateObjectIDFromHex(c.Id), Name: c.Name})
	}

	return coursesResult
}

func (c *GRPCServiceClient) Dial() (ps.CoursesServiceClient, error) {

	host := c.HOST + ":" + c.PORT

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not connect to %v %v", host, err))
	}

	log.Println("GRPC Connected to", host)

	c.Client = ps.NewCoursesServiceClient(conn)

	return c.Client, nil
}

func Construct(config contracts.Config) contracts.GRPCCourseService {
	return &GRPCServiceClient{HOST: config.GetAppConfig()["RPC_TARGET_HOST"], PORT: config.GetAppConfig()["RPC_TARGET_PORT"]}
}

//Procedural Test
//func serviceCourse() ps.CoursesServiceClient {
//
//	host := "172.24.0.3:6060"
//
//	conn, err := grpc.Dial(host, grpc.WithInsecure())
//	if err != nil {
//		log.Fatal("could not connect to", host, err)
//	}
//
//	return ps.NewCoursesServiceClient(conn)
//}
//
//func main() {
//
//	coursesServices := serviceCourse()
//
//	cID := ps.CoursesID{CoursesID: []string{"6300988647b1637e7974b3d9", "6300988647b1637e7974b3d6"}}
//
//	courses, err := coursesServices.List(context.TODO(), &cID)
//	if err != nil {
//		return
//	}
//
//	for _, c := range courses.List {
//		log.Println(c)
//	}
//
//}
