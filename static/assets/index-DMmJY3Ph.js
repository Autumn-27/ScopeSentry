import{cJ as n,cM as t,cN as r,bF as e,cO as u,cP as o,cQ as a,cg as i,cR as c,cS as s,bv as l,bD as f,cT as d,aV as v,bY as g}from"./index-vMt5tOuI.js";import{b}from"./isEqual-D2Iw95Gj.js";var p=1,h=2;function m(n){return n==n&&!t(n)}function j(n,t){return function(r){return null!=r&&(r[n]===t&&(void 0!==t||n in Object(r)))}}function y(t){var e=function(n){for(var t=r(n),e=t.length;e--;){var u=t[e],o=n[u];t[e]=[u,o,m(o)]}return t}(t);return 1==e.length&&e[0][2]?j(e[0][0],e[0][1]):function(r){return r===t||function(t,r,e,u){var o=e.length,a=o,i=!u;if(null==t)return!a;for(t=Object(t);o--;){var c=e[o];if(i&&c[2]?c[1]!==t[c[0]]:!(c[0]in t))return!1}for(;++o<a;){var s=(c=e[o])[0],l=t[s],f=c[1];if(i&&c[2]){if(void 0===l&&!(s in t))return!1}else{var d=new n;if(u)var v=u(l,f,s,t,r,d);if(!(void 0===v?b(f,l,p|h,u,d):v))return!1}}return!0}(r,t,e)}}function F(n,t){return null!=n&&t in Object(n)}function O(n,t){return null!=n&&function(n,t,r){for(var s=-1,l=(t=e(t,n)).length,f=!1;++s<l;){var d=u(t[s]);if(!(f=null!=n&&r(n,d)))break;n=n[d]}return f||++s!=l?f:!!(l=null==n?0:n.length)&&o(l)&&a(d,l)&&(i(n)||c(n))}(n,t,F)}var w=1,x=2;function E(n){return s(n)?(t=u(n),function(n){return null==n?void 0:n[t]}):function(n){return function(t){return f(t,n)}}(n);var t}function H(n){return"function"==typeof n?n:null==n?d:"object"==typeof n?i(n)?(t=n[0],r=n[1],s(t)&&m(r)?j(u(t),r):function(n){var e=l(n,t);return void 0===e&&e===r?O(n,t):b(r,e,w|x)}):y(n):E(n);var t,r}const M=new Map;let q;function A(n,t){let r=[];return Array.isArray(t.arg)?r=t.arg:g(t.arg)&&r.push(t.arg),function(e,u){const o=t.instance.popperRef,a=e.target,i=null==u?void 0:u.target,c=!t||!t.instance,s=!a||!i,l=n.contains(a)||n.contains(i),f=n===a,d=r.length&&r.some((n=>null==n?void 0:n.contains(a)))||r.length&&r.includes(i),v=o&&(o.contains(a)||o.contains(i));c||s||l||f||d||v||t.value(e,u)}}v&&(document.addEventListener("mousedown",(n=>q=n)),document.addEventListener("mouseup",(n=>{for(const t of M.values())for(const{documentHandler:r}of t)r(n,q)})));const L={beforeMount(n,t){M.has(n)||M.set(n,[]),M.get(n).push({documentHandler:A(n,t),bindingFn:t.value})},updated(n,t){M.has(n)||M.set(n,[]);const r=M.get(n),e=r.findIndex((n=>n.bindingFn===t.oldValue)),u={documentHandler:A(n,t),bindingFn:t.value};e>=0?r.splice(e,1,u):r.push(u)},unmounted(n){M.delete(n)}};export{L as C,H as b,O as h};