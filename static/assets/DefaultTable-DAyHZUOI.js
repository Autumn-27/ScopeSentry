import{_ as e}from"./ContentWrap.vue_vue_type_script_setup_true_lang-ByRQd7bd.js";import{d as t,l as a,v as o,e as l,C as i,r as s,o as r,i as m,w as p,a as n,I as d}from"./index-CN_-6CJN.js";import{_ as c}from"./Table.vue_vue_type_script_lang-B-A-w7yr.js";import{a as b}from"./index-BwFEEblY.js";import{E as j}from"./el-tag-C70eRmCH.js";import"./el-card-CNQMmRvy.js";import"./el-tooltip-l0sNRNKZ.js";import"./el-popper-qPvIts7N.js";import"./el-table-column-CwOTOpI2.js";import"./index-C-cnLFau.js";import"./debounce-CYJWxmBb.js";import"./el-checkbox-BWglfH5k.js";import"./isArrayLikeObject-B3dE5fSb.js";import"./raf-CKcl9dRI.js";import"./el-pagination-cWAg7nWs.js";import"./el-select-D1ImP6MX.js";import"./strings-BAothc3L.js";import"./useInput-BQndv-k3.js";import"./el-image-viewer-C6bUU2gm.js";import"./el-empty-Dch0KbwO.js";import"./tsxHelper-Bal41hqM.js";import"./el-dropdown-item-C_YpEPv1.js";import"./castArray-B2J6x547.js";import"./refs-BLlp70uH.js";import"./index-Cd5cNpU-.js";import"./index-CSisi_IG.js";const u=t({__name:"DefaultTable",setup(t){const{t:u}=a(),f=[{field:"title",label:u("tableDemo.title")},{field:"author",label:u("tableDemo.author")},{field:"display_time",label:u("tableDemo.displayTime"),sortable:!0},{field:"importance",label:u("tableDemo.importance"),formatter:(e,t,a)=>o(j,{type:1===a?"success":2===a?"warning":"danger"},(()=>u(1===a?"tableDemo.important":2===a?"tableDemo.good":"tableDemo.commonly")))},{field:"pageviews",label:u("tableDemo.pageviews")},{field:"action",label:u("tableDemo.action"),slots:{default:e=>{let t;return l(i,{type:"primary",onClick:()=>_(e)},"function"==typeof(a=t=u("tableDemo.action"))||"[object Object]"===Object.prototype.toString.call(a)&&!d(a)?t:{default:()=>[t]});var a}}}],g=s(!0);let y=s([]);(async e=>{const t=await b(e||{pageIndex:1,pageSize:10}).catch((()=>{})).finally((()=>{g.value=!1}));t&&(y.value=t.data.list)})();const _=e=>{};return(t,a)=>(r(),m(n(e),{title:n(u)("tableDemo.table"),message:n(u)("tableDemo.tableDes")},{default:p((()=>[l(n(c),{columns:f,data:n(y),loading:g.value,defaultSort:{prop:"display_time",order:"descending"}},null,8,["data","loading"])])),_:1},8,["title","message"]))}});export{u as default};