import{d as e,D as t,r as a,s as l,v as i,G as s,z as o,o as n,c as r,e as p,a as u,w as d,H as m,F as c,l as g,ai as j,_ as f}from"./index-ChGT_YCB.js";import{u as b}from"./useTable-DVIhlLGs.js";import{E as h}from"./el-card-Bv8eL-PY.js";import{E as v}from"./el-pagination-BHF_c2aK.js";import{E as y}from"./el-tag-uqONbW_Z.js";import"./el-select-BqJxfjvU.js";import"./el-popper-PhuldJD-.js";import{E as x,a as _}from"./el-col-i_K16CKa.js";import{_ as S}from"./Table.vue_vue_type_script_lang-BIM24Su6.js";import{u as E}from"./useCrudSchemas-BzAvyzCs.js";import{a as w,d as C,o as T}from"./index-C-GaKulA.js";import V from"./Csearch-BrOH0oMu.js";import"./index-CZTnDW2g.js";import"./strings-Dm6yfrSz.js";import"./useInput-D50pGW-V.js";import"./debounce-Zw0D3mea.js";import"./el-table-column-D3t4SgcR.js";import"./el-checkbox-CasWmTQM.js";import"./isArrayLikeObject-7_HOgyXV.js";import"./raf-CI9lY1B7.js";import"./el-tooltip-l0sNRNKZ.js";import"./el-image-viewer-DOJ4W1Wh.js";import"./el-empty-w9dzImmf.js";import"./tsxHelper-DC_KhjvT.js";import"./el-dropdown-item-DkQ0kmO2.js";import"./castArray-DD2tgwSj.js";import"./refs-Dfwo_yLk.js";import"./index-BAUEFq46.js";import"./tree-BfZhwLPs.js";import"./index-D1H-7IS5.js";import"./ContentWrap.vue_vue_type_script_setup_true_lang-JnnrwuOA.js";import"./el-text-iy-stWXi.js";import"./el-divider-Dm9ByY6A.js";import"./el-tree-select-B65aj2Zb.js";import"./index-BKDad1rC.js";/* empty css                */import"./el-switch-sUeo6Nig.js";import"./Dialog.vue_vue_type_style_index_0_lang-CjT1Och0.js";import"./useIcon-Cv7WwndZ.js";import"./exportData.vue_vue_type_script_setup_true_lang-BKfVLSbi.js";import"./el-tab-pane-DD3zcVo5.js";import"./el-form-D2bRQxI0.js";import"./el-radio-group-tIpiOwLu.js";import"./el-space-BVO1gQj3.js";const z=f(e({__name:"Subdomain",props:{projectList:{}},setup(e){const{t:f}=g(),z=[{keyword:"ip",example:'ip="192.168.2.1"',explain:f("searchHelp.ip")},{keyword:"domain",example:'domain="example.com"',explain:f("searchHelp.domain")},{keyword:"type",example:'type="CNAME"',explain:f("searchHelp.subdomainType")},{keyword:"value",example:'value="allcdn.example.com"',explain:f("searchHelp.subdoaminValue")},{keyword:"project",example:'project="Hackerone"',explain:f("searchHelp.project")}];t((()=>{H(),window.addEventListener("resize",H)}));const k=a(0),H=()=>{const e=window.innerHeight||document.documentElement.clientHeight;k.value=.7*e},A=a(""),N=e=>{A.value=e,B()},L=l({}),P=l([{field:"selection",type:"selection",minWidth:"55"},{field:"index",label:f("tableDemo.index"),type:"index",minWidth:"30"},{field:"host",label:f("subdomain.subdomainName"),minWidth:"200"},{field:"type",label:f("subdomain.recordType"),minWidth:"200",columnKey:"type",filters:[{text:"A",value:"A"},{text:"NS",value:"NS"},{text:"CNAME",value:"CNAME"},{text:"PTR",value:"PTR"},{text:"TXT",value:"TXT"}]},{field:"value",label:f("subdomain.recordValue"),minWidth:"250",formatter:(e,t,a)=>{let l="";return a.forEach(((e,t)=>{l+=`${e}\r\n`})),l}},{field:"ip",label:"IP",minWidth:"150",formatter:(e,t,a)=>{let l="";return a.forEach(((e,t)=>{l+=`${e}\r\n`})),l}},{field:"tags",label:"TAG",fit:"true",formatter:(e,t,l)=>{null==l&&(l=[]),L[e.id]||(L[e.id]={inputVisible:!1,inputValue:"",inputRef:a(null)});const n=L[e.id],r=async()=>{n.inputValue&&(l.push(n.inputValue),w(e.id,I,n.inputValue)),n.inputVisible=!1,n.inputValue=""};return i(_,{},(()=>[...l.map((t=>i(x,{span:24,key:t},(()=>[i("div",{onClick:e=>((e,t)=>{e.target.classList.contains("el-tag__close")||te("tags",t)})(e,t)},[i(y,{closable:!0,onClose:()=>(async t=>{const a=l.indexOf(t);a>-1&&l.splice(a,1),C(e.id,I,t)})(t)},(()=>t))])])))),i(x,{span:24},n.inputVisible?()=>i(s,{ref:n.inputRef,modelValue:n.inputValue,"onUpdate:modelValue":e=>n.inputValue=e,class:"w-20",size:"small",onKeyup:e=>{"Enter"===e.key&&r()},onBlur:r}):()=>i(o,{class:"button-new-tag",size:"small",onClick:()=>(n.inputVisible=!0,void j((()=>{})))},(()=>"+ New Tag")))]))},minWidth:"130"},{field:"time",label:f("asset.time"),minWidth:"200"}]);let I="subdomain";P.forEach((e=>{e.hidden=e.hidden??!1}));let W=a(!1);const U=({field:e,hidden:t})=>{const a=P.findIndex((t=>t.field===e));-1!==a&&(P[a].hidden=t),(()=>{const e=P.reduce(((e,t)=>(e[t.field]=t.hidden,e)),{});e.statisticsHidden=W.value,localStorage.setItem(`columnConfig_${I}`,JSON.stringify(e))})()};(()=>{const e=JSON.parse(localStorage.getItem(`columnConfig_${I}`)||"{}");P.forEach((t=>{void 0!==e[t.field]&&"select"!=t.field&&(t.hidden=e[t.field])})),W.value=e.statisticsHidden})();const{allSchemas:R}=E(P),{tableRegister:D,tableState:O,tableMethods:$}=b({fetchDataApi:async()=>{const{currentPage:e,pageSize:t}=O,a=await T(A.value,e.value,t.value,q);return{list:a.data.list,total:a.data.total}},immediate:!1}),{loading:F,dataList:M,total:K,currentPage:J,pageSize:X}=O,{getList:B,getElTableExpose:G}=$;function Y(){return{background:"var(--el-fill-color-light)"}}const q=l({}),Q=async e=>{Object.assign(q,e),B()},Z=(e,t)=>{Object.assign(q,t),A.value=e,B()},ee=a([]),te=(e,t)=>{const a=`${e}=${t}`;ee.value=[...ee.value,a]},ae=e=>{if(ee.value){const[t,a]=e.split("=");t in q&&Array.isArray(q[t])&&(q[t]=q[t].filter((e=>e!==a)),0===q[t].length&&delete q[t]),ee.value=ee.value.filter((t=>t!==e))}};return(e,t)=>(n(),r(c,null,[p(V,{getList:u(B),handleSearch:N,searchKeywordsData:z,index:"subdomain",getElTableExpose:u(G),projectList:e.$props.projectList,handleFilterSearch:Z,crudSchemas:P,dynamicTags:ee.value,handleClose:ae,onUpdateColumnVisibility:U},null,8,["getList","getElTableExpose","projectList","crudSchemas","dynamicTags"]),p(u(_),null,{default:d((()=>[p(u(x),null,{default:d((()=>[p(u(h),{style:{height:"min-content"}},{default:d((()=>[p(u(S),{pageSize:u(X),"onUpdate:pageSize":t[0]||(t[0]=e=>m(X)?X.value=e:null),currentPage:u(J),"onUpdate:currentPage":t[1]||(t[1]=e=>m(J)?J.value=e:null),columns:u(R).tableColumns,data:u(M),stripe:"",border:!0,loading:u(F),resizable:!0,onRegister:u(D),onFilterChange:Q,headerCellStyle:Y,style:{fontFamily:"-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji"}},null,8,["pageSize","currentPage","columns","data","loading","onRegister"])])),_:1})])),_:1}),p(u(x),{":span":24},{default:d((()=>[p(u(h),null,{default:d((()=>[p(u(v),{pageSize:u(X),"onUpdate:pageSize":t[2]||(t[2]=e=>m(X)?X.value=e:null),currentPage:u(J),"onUpdate:currentPage":t[3]||(t[3]=e=>m(J)?J.value=e:null),"page-sizes":[10,20,50,100,200,500,1e3],layout:"total, sizes, prev, pager, next, jumper",total:u(K)},null,8,["pageSize","currentPage","total"])])),_:1})])),_:1})])),_:1})],64))}}),[["__scopeId","data-v-5bb94ddb"]]);export{z as default};