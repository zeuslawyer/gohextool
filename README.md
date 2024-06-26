# Motivation

## Reference/Research

- [Function selectors](<https://medium.com/coinmonks/function-selectors-in-solidity-understanding-and-working-with-them-25e07755e976#:~:text=The%20function%20signature%20is%20derived,myFunction(address%2Cuint256)%20.>)

## Commands

1. Help : `hextool  --help << or -h>>` will list available commands. `hextool help <COMMAND> ` will print out the flags each command accepts or expects.

2. Hex to Int: `hextool toint --hex 0x0000000000000000000000000000000000000000000000000000000000690208` // 6881800

3. Hex to String `hextool tostring --hex 0x486578746f6f6c204d616b657320457468657265756d204465762045617369657221` // Hextool Makes Ethereum Dev Easier!

4. Retrieve function signature if given a function selector (<b>Note: </b> you must pass a path or url to a valid json object that has an `abi` property on it with an ABI array value). See below for examples.
   `hextool funcsig --selector 0xa9059cbb --url https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json` // transfer(address,uint256)

   or, using a path on your file system
   `hextool funcsig --path "/PATH/TO/FILE/erc20.abi.json" --selector 0x095ea7b3` // approve(address,uint256)

    <br>
    Please examine [the shape of the object](https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json) for this to work correctly. The ABI json files produced by Hardhat will work too.
    <br>

5. Calculate the function selector from the ABI-specified function signature (excludes the word 'function'): `hextool selector --sig 'transfer(address,uint256)'` // 0xa9059cbb
   <br>
   <b>Note: </b> The signature must be enclosed in single or double quotes.
   <br>

6 Decode a method's selector (4 bytes) into the function signature, given a valid ABI file. `hextool decodeMethodSelector --selector 0xa9059cbb --url "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json"`>

> should return transfer(address,uint256)

7 Decode an error's selector (4 bytes) into the error's signature, given a valid ABI file. `hextool decodeErrorSelector  --selector 0x07da6ee6 --url https://raw.githubusercontent.com/Cyfrin/ccip-contracts/main/contracts-ccip/abi/v0.8/Router.json`

> should return InsufficientFeeTokenAmount()

8. Decode the topic hash (on block explorers this shows up as `topic0`). The full 32-byte hexstring is needed, as is a valid ABI. `hextool decodeEvent --topic 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef --url "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json"`

9. Decode abi-encoded hex input with `hextool abi.decode --hex <<abi encoded hex>> --types "<<comma separated types in a string>>"`.  
   Example:

```
hextool abi.decode --hex "00000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000000060000000000000000000000000208aa722aca42399eac5192ee778e4d42f4e5de300000000000000000000000000000000000000000000000000000000000000057a7562696e000000000000000000000000000000000000000000000000000000" --types 'uint16,string,address'
```

Produces `[1001 zubin 0x208AA722Aca42399eaC5192EE778e4D42f4E5De3]`.

10. Abi-encode data (with their corresponding values). `hextool abi.encode --types '<<comma separated types in a string>>' --values '<<comma separated data string>>'`

eg: `hextool abi.encode --types 'uint64,string,address' --values '1981,bananas dude!,0x208AA722Aca42399eaC5192EE778e4D42f4E5De3'`
will produce: `0x00000000000000000000000000000000000000000000000000000000000007bd0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000208aa722aca42399eac5192ee778e4d42f4e5de3000000000000000000000000000000000000000000000000000000000000000d62616e616e617320647564652100000000000000000000000000000000000000`

Try and reverse that with `hextool abi.decode`!

## Other projects for research

https://github.com/umbracle/ethgo/tree/main //wrapper pkg
https://github.com/defiweb/go-eth // wrapper pkg

https://gist.github.com/crazygit/9279a3b26461d7cb03e807a6362ec855 // decoding tx logs, and reading contract ABI from etherscan
