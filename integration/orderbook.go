package integration

import (
	"encoding/json"
	"fmt"
	"github.com/tutengdihuang/okex_v5sdk_go/ws"
	"github.com/tutengdihuang/okex_v5sdk_go/ws/wImpl"
	"log"
	"strings"
	"time"
)

func prework() *ws.WsClient {
	//ep := "wss://wspap.okx.com:8443/ws/v5/public?brokerId=9999"
	ep := "wss://ws.okx.com:8443/ws/v5/public"
	r, err := ws.NewWsClient(ep)
	if err != nil {
		log.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		log.Fatal(err, ep)
	}
	return r
}

func OrderBooksRemote() {
	r := prework()
	var err error
	var res bool

	/*
		设置关闭深度数据管理
	*/
	// err = r.EnableAutoDepthMgr(false)
	// if err != nil {
	// 	fmt.Println("关闭自动校验失败！")
	// }

	end := make(chan struct{})
	var hookFunc = func(ts time.Time, data wImpl.DepthData) error {
		// 对于深度类型数据处理的用户可以自定义

		// 检测深度数据是否正常
		key, _ := json.Marshal(data.Arg)
		fmt.Println("个数：", len(data.Data[0].Asks))
		checksum := data.Data[0].Checksum
		fmt.Println("[自定义方法] ", string(key), ", checksum = ", checksum)

		for _, ask := range data.Data[0].Asks {

			arr := strings.Split(ask[0], ".")
			//fmt.Println(arr)
			if len(arr) > 1 && len(arr[1]) > 2 {
				fmt.Println("ask数据异常,", checksum, "ask:", ask)
				log.Fatal()
				end <- struct{}{}
				return nil
			} else {
				fmt.Println("bid数据正常,", checksum, "ask:", ask)
			}

		}

		for _, bid := range data.Data[0].Bids {

			arr := strings.Split(bid[0], ".")
			//fmt.Println(arr)
			if len(arr) > 1 && len(arr[1]) > 2 {
				fmt.Println("bid数据异常,", checksum, "bid:", bid)
				log.Fatal()
				end <- struct{}{}
				return nil
			} else {
				fmt.Println("ask数据正常,", checksum, "bid:", bid)
			}

		}

		// 查看当前合并后的全量深度数据
		snapshot, err := r.GetSnapshotByChannel(data)
		if err != nil {
			log.Fatal("深度数据不存在！")
		}
		// 展示ask/bid 前5档数据
		fmt.Println(" Ask 5 档数据 >> ")
		for _, v := range snapshot.Asks[:5] {
			fmt.Println(" price:", v[0], " amount:", v[1])
		}
		fmt.Println(" Bid 5 档数据 >> ")
		for _, v := range snapshot.Bids[:5] {
			fmt.Println(" price:", v[0], " amount:", v[1])
		}
		return nil
	}
	r.AddDepthHook(hookFunc)

	// 可选类型：books books5 books-l2-tbt
	channel := "books-l2-tbt"

	instIds := []string{"BTC-USDT"}
	for _, instId := range instIds {
		var args []map[string]string
		arg := make(map[string]string)
		arg["instId"] = instId
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PubOrderBooks(ws.OP_SUBSCRIBE, channel, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("订阅成功！", usedTime.String())
		} else {
			fmt.Println("订阅失败！", err)
			log.Fatal("订阅失败！", err)
		}
	}

	select {
	case <-end:

	}
	//等待推送
	for _, instId := range instIds {
		var args []map[string]string
		arg := make(map[string]string)
		arg["instId"] = instId
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PubOrderBooks(ws.OP_UNSUBSCRIBE, channel, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("取消订阅成功！", usedTime.String())
		} else {
			fmt.Println("取消订阅失败！", err)
			log.Fatal("取消订阅失败！", err)
		}
	}

}

func OrderBooksRemoteSlot1() {
	r := prework()
	var err error
	var res bool

	/*
		设置关闭深度数据管理
	*/
	// err = r.EnableAutoDepthMgr(false)
	// if err != nil {
	// 	fmt.Println("关闭自动校验失败！")
	// }

	end := make(chan struct{})
	var hookFunc = func(ts time.Time, data wImpl.MsgData) error {
		// 对于深度类型数据处理的用户可以自定义

		// 检测深度数据是否正常
		fmt.Println("个数：", len(data.Data))
		fmt.Printf("===%+v\n", data)
		fmt.Println("===", data)
		return nil
	}
	r.AddBookMsgHook(hookFunc)

	// 可选类型：books books5 books-l2-tbt
	channel := "bbo-tbt"

	instIds := []string{"ETH-USDT"}
	for _, instId := range instIds {
		var args []map[string]string
		arg := make(map[string]string)
		arg["instId"] = instId
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PubOrderBooks(ws.OP_SUBSCRIBE, channel, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("订阅成功！", usedTime.String())
		} else {
			fmt.Println("订阅失败！", err)
			log.Fatal("订阅失败！", err)
		}
	}

	select {
	case <-end:

	}
	//等待推送
	for _, instId := range instIds {
		var args []map[string]string
		arg := make(map[string]string)
		arg["instId"] = instId
		args = append(args, arg)

		start := time.Now()
		res, _, err = r.PubOrderBooks(ws.OP_UNSUBSCRIBE, channel, args)
		if res {
			usedTime := time.Since(start)
			fmt.Println("取消订阅成功！", usedTime.String())
		} else {
			fmt.Println("取消订阅失败！", err)
			log.Fatal("取消订阅失败！", err)
		}
	}

}
