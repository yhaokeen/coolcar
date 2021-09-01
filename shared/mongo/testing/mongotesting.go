package mongotesting

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"testing"
)

const (
	image         = "mongo:4.4"
	containerPort = "27017/tcp"
)

func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	//创建容器
	containerBody, err := c.ContainerCreate(ctx, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			containerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	//删除容器
	defer func() {
		err := c.ContainerRemove(ctx, containerBody.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()
	//启动容器
	err = c.ContainerStart(ctx, containerBody.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	//获取容器的配置
	inspRsp, err := c.ContainerInspect(ctx, containerBody.ID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRsp.NetworkSettings.Ports["27017/tcp"][0]
	*mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}
