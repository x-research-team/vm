## Run vm programs in the browser using WebAssembly

In the vm home directory, just run `run_wasm.sh`, then
open your browser, and type `http://localhost:9090`.


Below is the snapshot of the running vm demo in browser:

<p>
    <img alt="vm playground" src="https://github.com/haifenghuang/vm/blob/master/wasm/vm_playground.png?raw=true" width="450" height="450">
</p>

<br>
<p>
    <img alt="vm playground" src="https://github.com/haifenghuang/vm/blob/master/wasm/vm_playground2.png?raw=true" width="450" height="450">
</p>

<br>
<p>
    <img alt="vm playground" src="https://github.com/haifenghuang/vm/blob/master/wasm/vm_playground3.png?raw=true" width="450" height="450">
</p>

## Limitation

1. Can not use 'stdin', 'stdout' and 'stderr'
2. Can not use file object's method. e.g. fileObj.read(), fileObj.readLine()