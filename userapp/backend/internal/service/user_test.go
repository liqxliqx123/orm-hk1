package service

// func TestNewUserService(t *testing.T) {
// 	c := buildCache(t)
// 	s := NewUserService(c, nil)
//
// 	ctx := context.Background()
// 	reply, err := s.CreateUser(ctx, &dto.CreateUserReq{
// 		User: &dto.User{
// 			Email: fmt.Sprintf("%d@xxmail.com", time.Now().Unix()),
// 			Password: "abc",
// 		},
// 	})
// 	assert.Nil(t, err)
// 	id := reply.UserId
// 	assert.True(t, id > 0)
//
// 	get, err := s.FindById(ctx, &dto.FindByUserIdReq{
// 		Id: id,
// 	})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, get)
//
// 	get, err = s.FindById(ctx, &dto.FindByUserIdReq{
// 		Id: id,
// 	})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, get)
// }

// func buildCache(t *testing.T) cache.Cache {
// 	rc := redis.NewClient(&redis.Options{
// 		Addr: "localhost:6379",
// 		Password: "abc",
// 	})
// 	c, err := cache.NewRedisCache(rc, cache.WithMarshalFunc(func(m interface{}) ([]byte, error) {
// 			return proto.Marshal(m.(proto.Message))
// 		}), cache.WithUnmarshalFunc(func(b []byte, m interface{}) error {
// 			return proto.Unmarshal(b, m.(proto.Message))
// 		}))
//
// 	assert.Nil(t, err)
// 	return c
// }
