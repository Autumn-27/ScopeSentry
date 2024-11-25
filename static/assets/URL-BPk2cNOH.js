import{d as e,r as t,s as l,e as a,y as i,S as s,v as o,G as r,z as n,D as p,o as u,c as d,a as c,w as m,H as g,F as f,I as j,l as v,ai as h,_ as x}from"./index-CN_-6CJN.js";import{u as y}from"./useTable-D4uu-Soo.js";import{E as b}from"./el-card-CNQMmRvy.js";import{E as _}from"./el-pagination-cWAg7nWs.js";import{E as S}from"./el-tag-C70eRmCH.js";import"./el-select-D1ImP6MX.js";import"./el-popper-qPvIts7N.js";import{E,a as w}from"./el-col-jGixOe2y.js";import{E as L}from"./el-text-Bb8H_RHf.js";import{_ as z}from"./Table.vue_vue_type_script_lang-B-A-w7yr.js";import{u as C}from"./useCrudSchemas-Di-bMQjR.js";import{a as U,d as k,q as V}from"./index-BiRfkCSs.js";import H from"./Csearch-Dagrp5kp.js";import"./index-C-cnLFau.js";import"./strings-BAothc3L.js";import"./useInput-BQndv-k3.js";import"./debounce-CYJWxmBb.js";import"./el-table-column-CwOTOpI2.js";import"./el-checkbox-BWglfH5k.js";import"./isArrayLikeObject-B3dE5fSb.js";import"./raf-CKcl9dRI.js";import"./el-tooltip-l0sNRNKZ.js";import"./el-image-viewer-C6bUU2gm.js";import"./el-empty-Dch0KbwO.js";import"./tsxHelper-Bal41hqM.js";import"./el-dropdown-item-C_YpEPv1.js";import"./castArray-B2J6x547.js";import"./refs-BLlp70uH.js";import"./index-Cd5cNpU-.js";import"./tree-BfZhwLPs.js";import"./index-CSisi_IG.js";import"./ContentWrap.vue_vue_type_script_setup_true_lang-ByRQd7bd.js";import"./el-divider-UM5pYHfY.js";import"./el-tree-select-CLQSZiq_.js";import"./index-BakKIM4W.js";import"./el-switch-C_vE8S4c.js";import"./Dialog.vue_vue_type_style_index_0_lang-gBBmaUwO.js";import"./useIcon-CoAv5gAq.js";import"./exportData.vue_vue_type_script_setup_true_lang-BD7aTcXR.js";import"./el-tab-pane-Cb5dGpe_.js";import"./el-form-gAlh6c2U.js";import"./el-radio-group-cVqelweb.js";import"./el-space-vXK55lKA.js";const W=x(e({__name:"URL",props:{projectList:{}},setup(e){const{t:x}=v(),W=[{keyword:"url",example:'url="http://example.com"',explain:x("searchHelp.url")},{keyword:"input",example:'input="example.com"',explain:x("searchHelp.inpur")},{keyword:"source",example:'source="exapmle.com/example.js"',explain:x("searchHelp.source")},{keyword:"type",example:'type="linkfinder"',explain:x("searchHelp.urlType")},{keyword:"project",example:'project="Hackerone"',explain:x("searchHelp.project")}],R=t(""),I=e=>{R.value=e,X()},T=l({}),A=l([{field:"selection",type:"selection",minWidth:"55"},{field:"index",label:x("tableDemo.index"),type:"index",minWidth:55},{field:"url",label:"URL",minWidth:250},{field:"status",label:x("dirScan.status"),columnKey:"status",minWidth:120,formatter:(e,t,l)=>{if(null==l)return a("div",null,[i("-")]);let o="";return o=l<300?"#2eb98a":"#ff5252",a(w,{gutter:20},{default:()=>[a(E,{span:1},{default:()=>[a(s,{icon:"clarity:circle-solid",color:o,size:10,style:"transform: translateY(8%)"},null)]}),a(E,{span:2},{default:()=>{return[a(L,null,(e=l,"function"==typeof e||"[object Object]"===Object.prototype.toString.call(e)&&!j(e)?l:{default:()=>[l]}))];var e}})]})},filters:[{text:"200",value:200},{text:"201",value:201},{text:"204",value:204},{text:"301",value:301},{text:"302",value:302},{text:"304",value:304},{text:"400",value:400},{text:"401",value:401},{text:"403",value:403},{text:"404",value:404},{text:"500",value:500},{text:"502",value:502},{text:"503",value:503},{text:"504",value:504}]},{field:"length",label:"Length",minWidth:120,sortable:"custom"},{field:"source",label:x("URL.source"),minWidth:100},{field:"type",label:x("URL.type"),minWidth:100},{field:"input",label:x("URL.input"),minWidth:200},{field:"tags",label:"TAG",fit:"true",formatter:(e,l,a)=>{null==a&&(a=[]),T[e.id]||(T[e.id]={inputVisible:!1,inputValue:"",inputRef:t(null)});const i=T[e.id],s=async()=>{i.inputValue&&(a.push(i.inputValue),U(e.id,O,i.inputValue)),i.inputVisible=!1,i.inputValue=""};return o(w,{},(()=>[...a.map((t=>o(E,{span:24,key:t},(()=>[o("div",{onClick:e=>((e,t)=>{e.target.classList.contains("el-tag__close")||re("tags",t)})(e,t)},[o(S,{closable:!0,onClose:()=>(async t=>{const l=a.indexOf(t);l>-1&&a.splice(l,1),k(e.id,O,t)})(t)},(()=>t))])])))),o(E,{span:24},i.inputVisible?()=>o(r,{ref:i.inputRef,modelValue:i.inputValue,"onUpdate:modelValue":e=>i.inputValue=e,class:"w-20",size:"small",onKeyup:e=>{"Enter"===e.key&&s()},onBlur:s}):()=>o(n,{class:"button-new-tag",size:"small",onClick:()=>(i.inputVisible=!0,void h((()=>{})))},(()=>"+ New Tag")))]))},minWidth:"130"},{field:"time",label:x("asset.time"),minWidth:200}]);let O="UrlScan";A.forEach((e=>{e.hidden=e.hidden??!1}));let P=t(!1);const D=({field:e,hidden:t})=>{const l=A.findIndex((t=>t.field===e));-1!==l&&(A[l].hidden=t),(()=>{const e=A.reduce(((e,t)=>(e[t.field]=t.hidden,e)),{});e.statisticsHidden=P.value,localStorage.setItem(`columnConfig_${O}`,JSON.stringify(e))})()};(()=>{const e=JSON.parse(localStorage.getItem(`columnConfig_${O}`)||"{}");A.forEach((t=>{void 0!==e[t.field]&&"select"!=t.field&&(t.hidden=e[t.field])})),P.value=e.statisticsHidden})();const{allSchemas:N}=C(A),{tableRegister:F,tableState:$,tableMethods:K}=y({fetchDataApi:async()=>{const{currentPage:e,pageSize:t}=$,l=await V(R.value,e.value,t.value,ae,Q);return{list:l.data.list,total:l.data.total}},immediate:!1}),{loading:J,dataList:M,total:B,currentPage:G,pageSize:q}=$,{getList:X,getElTableExpose:Y}=K;function Z(){return{background:"var(--el-fill-color-light)"}}q.value=20,p((()=>{le(),window.addEventListener("resize",le)}));const Q=l({}),ee=async e=>{const t=e.prop,l=e.order;Q[t]=l,X()},te=t(0),le=()=>{const e=window.innerHeight||document.documentElement.clientHeight;te.value=.7*e},ae=l({}),ie=(e,t)=>{Object.assign(ae,t),R.value=e,X()},se=async e=>{Object.assign(ae,e),X()},oe=t([]),re=(e,t)=>{const l=`${e}=${t}`;oe.value=[...oe.value,l]},ne=e=>{if(oe.value){const[t,l]=e.split("=");t in ae&&Array.isArray(ae[t])&&(ae[t]=ae[t].filter((e=>e!==l)),0===ae[t].length&&delete ae[t]),oe.value=oe.value.filter((t=>t!==e))}};return(e,t)=>(u(),d(f,null,[a(H,{getList:c(X),handleSearch:I,searchKeywordsData:W,index:"UrlScan",getElTableExpose:c(Y),projectList:e.$props.projectList,handleFilterSearch:ie,crudSchemas:A,dynamicTags:oe.value,handleClose:ne,onUpdateColumnVisibility:D},null,8,["getList","getElTableExpose","projectList","crudSchemas","dynamicTags"]),a(c(w),null,{default:m((()=>[a(c(E),null,{default:m((()=>[a(c(b),null,{default:m((()=>[a(c(z),{pageSize:c(q),"onUpdate:pageSize":t[0]||(t[0]=e=>g(q)?q.value=e:null),currentPage:c(G),"onUpdate:currentPage":t[1]||(t[1]=e=>g(G)?G.value=e:null),columns:c(N).tableColumns,data:c(M),"max-height":te.value,stripe:"",border:!0,loading:c(J),resizable:!0,onSortChange:ee,onRegister:c(F),onFilterChange:se,headerCellStyle:Z,style:{fontFamily:"-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji"}},null,8,["pageSize","currentPage","columns","data","max-height","loading","onRegister"])])),_:1})])),_:1}),a(c(E),{":span":24},{default:m((()=>[a(c(b),null,{default:m((()=>[a(c(_),{pageSize:c(q),"onUpdate:pageSize":t[2]||(t[2]=e=>g(q)?q.value=e:null),currentPage:c(G),"onUpdate:currentPage":t[3]||(t[3]=e=>g(G)?G.value=e:null),"page-sizes":[20,50,100,200,500,1e3],layout:"total, sizes, prev, pager, next, jumper",total:c(B)},null,8,["pageSize","currentPage","total"])])),_:1})])),_:1})])),_:1})],64))}}),[["__scopeId","data-v-48fb44f2"]]);export{W as default};