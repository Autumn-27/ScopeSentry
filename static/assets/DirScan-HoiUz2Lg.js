import{d as e,r as t,s as a,e as l,y as i,S as s,v as o,G as r,z as n,D as p,o as u,c as d,a as c,w as m,H as g,F as f,I as j,l as h,ai as v,_ as x}from"./index-lnK3vi2u.js";import{u as y}from"./useTable-BjrIXiJu.js";import{E as b}from"./el-card-Bnv5_xAI.js";import{E as _}from"./el-pagination-B5C7bQ26.js";import{E as S}from"./el-tag-BuYucu0j.js";import"./el-select-S7B-i-tn.js";import"./el-popper-u58bCfVz.js";import{E as w,a as E}from"./el-col-BaaUI79S.js";import{E as C}from"./el-text-B59PEWrW.js";import{_ as z}from"./Table.vue_vue_type_script_lang-Ddm6485n.js";import{u as k}from"./useCrudSchemas-D-VS2E3G.js";import{a as V,d as H,i as L}from"./index-D4vyQVMi.js";import A from"./Csearch-n5lZERfl.js";import"./index-W8_TjnTK.js";import"./strings-D_ONF1MI.js";import"./useInput-C2svBfmf.js";import"./debounce-D1p90kO3.js";import"./el-table-column-0AlIaPKN.js";import"./el-checkbox-DZ6O993m.js";import"./isArrayLikeObject-CXclMpK0.js";import"./raf-BAxP01Hz.js";import"./el-tooltip-l0sNRNKZ.js";import"./el-image-viewer-CJ6msFhC.js";import"./el-empty-Cpo3UJLh.js";import"./tsxHelper-Dd9KbMFH.js";import"./el-dropdown-item-5WGALSI3.js";import"./castArray-Vw7lTMUi.js";import"./refs-CYfXX6dK.js";import"./index-DaGTYmBc.js";import"./tree-BfZhwLPs.js";import"./index-BlMFUDDZ.js";import"./ContentWrap.vue_vue_type_script_setup_true_lang-iVEUw-M_.js";import"./el-divider-DFBHi6tM.js";import"./el-tree-select-DWBLoOYU.js";import"./index-CKqeiIRX.js";import"./el-switch-CNhl8tfg.js";import"./Dialog.vue_vue_type_style_index_0_lang-t63tDEPk.js";import"./useIcon-Dl780yOf.js";import"./exportData.vue_vue_type_script_setup_true_lang-B-9oyrSm.js";import"./el-tab-pane-Ck8OoNHB.js";import"./el-form-BQdiY2FK.js";import"./el-radio-group-D9C0bcyd.js";import"./el-space-B4SifBdK.js";const U=x(e({__name:"DirScan",props:{projectList:{}},setup(e){const{t:x}=h(),U=[{keyword:"url",example:'url="http://example.com"',explain:x("searchHelp.url")},{keyword:"statuscode",example:'statuscode=="200"',explain:x("searchHelp.statuscode")},{keyword:"redirect",example:'redirect="https://example.com"',explain:x("searchHelp.redirect")},{keyword:"project",example:'project="Hackerone"',explain:x("searchHelp.project")},{keyword:"length",example:'length="1234"',explain:x("searchHelp.length")}],D=t(""),I=e=>{D.value=e,Y()},R=a({}),O=a([{field:"selection",type:"selection",minWidth:"55"},{field:"index",label:x("tableDemo.index"),type:"index",minWidth:55},{field:"url",label:"URL",minWidth:200},{field:"status",label:x("dirScan.status"),columnKey:"status",minWidth:120,formatter:(e,t,a)=>{if(null==a)return l("div",null,[i("-")]);let o="";return o=a<300?"#2eb98a":"#ff5252",l(E,{gutter:1},{default:()=>[l(w,{span:1},{default:()=>[l(s,{icon:"clarity:circle-solid",color:o,size:10,style:"transform: translateY(8%)"},null)]}),l(w,{span:2},{default:()=>{return[l(C,null,(e=a,"function"==typeof e||"[object Object]"===Object.prototype.toString.call(e)&&!j(e)?a:{default:()=>[a]}))];var e}})]})},filters:[{text:"200",value:200},{text:"201",value:201},{text:"204",value:204},{text:"301",value:301},{text:"302",value:302},{text:"304",value:304},{text:"400",value:400},{text:"401",value:401},{text:"403",value:403},{text:"404",value:404},{text:"500",value:500},{text:"502",value:502},{text:"503",value:503},{text:"504",value:504}]},{field:"length",label:"Length",minWidth:120,sortable:"custom"},{field:"msg",label:"Redirect",minWidth:200},{field:"tags",label:"TAG",fit:"true",formatter:(e,a,l)=>{null==l&&(l=[]),R[e.id]||(R[e.id]={inputVisible:!1,inputValue:"",inputRef:t(null)});const i=R[e.id],s=async()=>{i.inputValue&&(l.push(i.inputValue),V(e.id,T,i.inputValue)),i.inputVisible=!1,i.inputValue=""};return o(E,{},(()=>[...l.map((t=>o(w,{span:24,key:t},(()=>[o("div",{onClick:e=>((e,t)=>{e.target.classList.contains("el-tag__close")||re("tags",t)})(e,t)},[o(S,{closable:!0,onClose:()=>(async t=>{const a=l.indexOf(t);a>-1&&l.splice(a,1),H(e.id,T,t)})(t)},(()=>t))])])))),o(w,{span:24},i.inputVisible?()=>o(r,{ref:i.inputRef,modelValue:i.inputValue,"onUpdate:modelValue":e=>i.inputValue=e,class:"w-20",size:"small",onKeyup:e=>{"Enter"===e.key&&s()},onBlur:s}):()=>o(n,{class:"button-new-tag",size:"small",onClick:()=>(i.inputVisible=!0,void v((()=>{})))},(()=>"+ New Tag")))]))},minWidth:"130"}]);let T="DirScanResult";O.forEach((e=>{e.hidden=e.hidden??!1}));let W=t(!1);const P=({field:e,hidden:t})=>{const a=O.findIndex((t=>t.field===e));-1!==a&&(O[a].hidden=t),(()=>{const e=O.reduce(((e,t)=>(e[t.field]=t.hidden,e)),{});e.statisticsHidden=W.value,localStorage.setItem(`columnConfig_${T}`,JSON.stringify(e))})()};(()=>{const e=JSON.parse(localStorage.getItem(`columnConfig_${T}`)||"{}");O.forEach((t=>{void 0!==e[t.field]&&"select"!=t.field&&(t.hidden=e[t.field])})),W.value=e.statisticsHidden})();const{allSchemas:N}=k(O),{tableRegister:F,tableState:$,tableMethods:K}=y({fetchDataApi:async()=>{const{currentPage:e,pageSize:t}=$,a=await L(D.value,e.value,t.value,te,le);return{list:a.data.list,total:a.data.total}},immediate:!1}),{loading:J,dataList:M,total:B,currentPage:G,pageSize:X}=$,{getList:Y,getElTableExpose:Z}=K;function q(){return{background:"var(--el-fill-color-light)"}}X.value=20,p((()=>{ee(),window.addEventListener("resize",ee)}));const Q=t(0),ee=()=>{const e=window.innerHeight||document.documentElement.clientHeight;Q.value=.7*e},te=a({}),ae=async e=>{Object.assign(te,e),Y()},le=a({}),ie=async e=>{const t=e.prop,a=e.order;le[t]=a,Y()},se=(e,t)=>{Object.assign(te,t),D.value=e,Y()},oe=t([]),re=(e,t)=>{const a=`${e}=${t}`;oe.value=[...oe.value,a]},ne=e=>{if(oe.value){const[t,a]=e.split("=");t in te&&Array.isArray(te[t])&&(te[t]=te[t].filter((e=>e!==a)),0===te[t].length&&delete te[t]),oe.value=oe.value.filter((t=>t!==e))}};return(e,t)=>(u(),d(f,null,[l(A,{getList:c(Y),handleSearch:I,searchKeywordsData:U,index:"DirScanResult",dynamicTags:oe.value,handleClose:ne,getElTableExpose:c(Z),handleFilterSearch:se,projectList:e.$props.projectList,crudSchemas:O,onUpdateColumnVisibility:P},null,8,["getList","dynamicTags","getElTableExpose","projectList","crudSchemas"]),l(c(E),null,{default:m((()=>[l(c(w),null,{default:m((()=>[l(c(b),null,{default:m((()=>[l(c(z),{pageSize:c(X),"onUpdate:pageSize":t[0]||(t[0]=e=>g(X)?X.value=e:null),currentPage:c(G),"onUpdate:currentPage":t[1]||(t[1]=e=>g(G)?G.value=e:null),columns:c(N).tableColumns,data:c(M),stripe:"",border:!0,loading:c(J),resizable:!0,"max-height":Q.value,onRegister:c(F),onFilterChange:ae,onSortChange:ie,headerCellStyle:q,"tooltip-options":{offset:1,showArrow:!1,effect:"dark",enterable:!0,showAfter:0,popperOptions:{},popperClass:"test",placement:"top",hideAfter:0,disabled:!0},style:{fontFamily:"-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji"}},null,8,["pageSize","currentPage","columns","data","loading","max-height","onRegister"])])),_:1})])),_:1}),l(c(w),{":span":24},{default:m((()=>[l(c(b),null,{default:m((()=>[l(c(_),{pageSize:c(X),"onUpdate:pageSize":t[2]||(t[2]=e=>g(X)?X.value=e:null),currentPage:c(G),"onUpdate:currentPage":t[3]||(t[3]=e=>g(G)?G.value=e:null),"page-sizes":[10,20,50,100,200,500,1e3],layout:"total, sizes, prev, pager, next, jumper",total:c(B)},null,8,["pageSize","currentPage","total"])])),_:1})])),_:1})])),_:1})],64))}}),[["__scopeId","data-v-ef90a4e7"]]);export{U as default};