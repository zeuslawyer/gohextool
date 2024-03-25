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

6 Decode a selector (4 bytes) into the function signature, given a valid ABI file. `hextool decodeSelector --selector 0xa9059cbb --url "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json"`

7. Decode the topic hash (on block explorers this shows up as `topic0`). The full 32-byte hexstring is needed, as is a valid ABI. `hextool decodeEvent --topic 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef --url "https://gist.githubusercontent.com/zeuslawyer/ecec03ff3f50311e510c201de4c076d5/raw/f096531942e922cb3f1d5daa2132f0e476356ced/good-data-erc20.json"`
