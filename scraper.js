let p = new DOMParser()

const list = []

let activeRequests = 0

const fetcher = (codepoint) => {
    activeRequests++
    setTimeout(() => {
        fetch(`https://codepoints.net/U+${codepoint}?lang=en`)
            .then(res => res.text())
        .then(res => {
            console.log('resource fetched for codepoint',codepoint)
            let x = p.parseFromString(res,'text/html')
            list.push({unicode:"U+FFFD", utf8: x.querySelector('tr[data-system=utf-8] > td').textContent.split(' ').map(it => parseInt(`0x${it}`)), utf16: x.querySelector('tr[data-system=utf-16] > td').textContent.split(' ').map(it => parseInt(`0x${it}`)), utf32: x.querySelector('tr[data-system=utf-32] > td').textContent.split(' ').map(it => parseInt(`0x${it}`))})
        })
        .catch(err => {
            console.log('failed to fetch resource for codepoint',codepoint,err)
        })
        .finally(() => {
            if(activeRequests > 0) {
                activeRequests--;
            }
        })
    },activeRequests*10)
}

let hex = ['0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f']
let visited = {}
let times = 0
function recurse(i, j, k, l){
    if(visited[hex[i]+hex[j]+hex[k]+hex[l]]) return
    // console.log([hex[i],hex[j], hex[k], hex[l]])
    visited[hex[i]+hex[j]+hex[k]+hex[l]] = true
    times++
    fetcher(hex[i]+hex[j]+hex[k]+hex[l])
    if(i+1 < hex.length) {
        recurse(i+1, j, k, l)
    }
    if(j+1 < hex.length) {
        recurse(i, j+1, k, l)
    }
    if(k+1 < hex.length){
        recurse(i, j, k+1, l)
    }
    if(l+1 < hex.length){
        recurse(i, j, k, l+1)
    }
}
recurse(0, 0, 0, 0)
console.log('times',times)