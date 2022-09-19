package pair

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/DonggyuLim/uniswap/client"
	"github.com/DonggyuLim/uniswap/db"
	"github.com/DonggyuLim/uniswap/pool"
	"github.com/DonggyuLim/uniswap/utils"
)

type Pair struct {
	Name string    `json:"name"`
	Pool pool.Pool `json:"pool"`
}

func CreatePair(tokenA, tokenB pool.Token) (Pair, error) {
	//erc20 에 grpc 로 통신하여 tokenA 와 tokenB 가 있는지 확인해야함.
	err := checkSameToken(tokenA, tokenB)
	if err != nil {
		return Pair{}, err
	}
	//만약 usdt 와 dai 페어를 만든다면 어떻게 해야할까?
	//map[usdt] = dai
	//map[dai] = usdt
	//이런식으로 만들기보단 무조건 문자열 순으로 들어오게 해주면 어떨가?
	/*
	   uniswap v2
	   getPair[token0][token1] = pair;
	   getPair[token1][token0] = pair; // populate mapping in the reverse direction
	   allPairs.push(pair);
	*/
	//유니스왑에서는 2차원 맵으로 처리를 해버린다.
	//이렇게하면 오히려 더 메모리를 잡아먹지 않을까?
	//tokenA 와 tokenB문자열을 합쳐버려서 보관하면 어떨까?
	//"usdt:dai","dai:usdt" 이런식으로 pair 를 지정해버리면 되지 않을까?
	// 아니면 tokenA 와 tokenB 의 문자열을 sort 해서 키로 만들어주면 될 것 같다.
	key := GetKey(tokenA.Name, tokenB.Name)
	pc := pool.NewPoolToken(key, getSymbol(tokenA.Name, tokenB.Name), 10)
	p := Pair{
		Name: key,
		Pool: pool.CreatePool(tokenA, tokenB, pc),
	}
	b, err := client.GetClient().GetBalance("USDT", "0xb")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("balance----------------", b)
	if err != nil {
		return Pair{}, err
	}
	err = db.Add("pair", p.GetName(), utils.DataToByte(p))
	if err != nil {
		return Pair{}, err
	}
	return p, nil
}

func (p *Pair) GetName() string {
	return p.Name
}

func checkSameToken(tokenA, tokenB pool.Token) error {
	if tokenA.Name == tokenB.Name {
		return errors.New("same tokens can't become pair")
	}
	return nil
}

// usdt,dai => dai:usdt
func GetKey(a, b string) string {
	slice := []string{a, b}
	sort.Strings(slice)
	return fmt.Sprintf("%s:%s", slice[0], slice[1])
}

// usdt,dai => dnuLP
func getSymbol(a, b string) string {
	aslice := strings.Split(a, "")
	bslice := strings.Split(b, "")
	result := fmt.Sprintf("%sN%sLP", aslice[0], bslice[0])
	return result
}

func ByteToPair(data []byte) (Pair, error) {
	var pair Pair
	decoder := gob.NewDecoder(bytes.NewBuffer((data)))
	err := decoder.Decode(&pair)
	if err != nil {
		return pair, err
	}
	return pair, nil
}
