import{_ as e}from"./ContentWrap.vue_vue_type_script_setup_true_lang-F13nrP2r.js";import{d as t,r as a,s as o,v as s,G as l,e as i,F as n,H as r,o as u,c as d,w as c,a as p,y as m,t as g,J as f,z as v,f as y,K as h,L as k,l as _}from"./index-KWTKF9N9.js";import{a as b,E as S}from"./el-col-t_F-vIr3.js";import{E as w}from"./el-text-bFj0OZ8x.js";import{E as j}from"./el-progress-FyGDfwQU.js";import{E as x}from"./el-tag-MzjVikBY.js";import{_ as T}from"./Table.vue_vue_type_script_lang-iUd8q4UT.js";import{u as C}from"./useTable-CuwOU5m7.js";import{u as I}from"./useIcon-Ry17CHl2.js";import{a as P,d as N,r as A,g as E}from"./index-8Y4g3RH2.js";import{_ as L}from"./Dialog.vue_vue_type_style_index_0_lang-T5U2It3I.js";import{_ as z}from"./AddTask.vue_vue_type_script_setup_true_lang-RCWLy57l.js";import{_ as U}from"./ProgressInfo.vue_vue_type_script_setup_true_lang-7qe38ALB.js";import"./el-card-9F69xyLv.js";import"./el-tooltip-w40geAFS.js";import"./el-popper-pnFa2Ecw.js";import"./el-checkbox-U1iRgEY7.js";import"./useInput-CGBP5_dt.js";import"./debounce-QvfGQs-0.js";import"./el-pagination-ooblyvhx.js";import"./el-image-viewer-S1bd3VyN.js";import"./tsxHelper-NOZkKkVH.js";import"./el-dropdown-item-bjCsATwt.js";import"./castArray-9Ouy5Ghy.js";import"./refs-WVQk01D0.js";import"./index-BS1c_xid.js";import"./raf-kE1zktlZ.js";import"./index-XFn6QxJu.js";/* empty css                          */import"./el-divider-685xvbVJ.js";import"./el-form-item-i8ZQhvu8.js";import"./el-switch-tqt83myp.js";import"./el-radio-group-MC2AFabc.js";import"./el-select-v2-0dTEd5GK.js";import"./el-input-number-Rb_QWPs-.js";import"./index-Dm_8ErkC.js";import"./index-kUGlyjnD.js";import"./index-6VLehrwJ.js";const V={class:"mb-10px"},D={style:{position:"relative",top:"12px"}};function F(e){return"function"==typeof e||"[object Object]"===Object.prototype.toString.call(e)&&!k(e)}const M=t({__name:"Task",setup(t){const k=I({icon:"iconoir:search"}),{t:M}=_(),W=a(""),H=()=>{te()},R=o([{field:"selection",type:"selection",width:"55"},{field:"name",label:M("task.taskName"),minWidth:30},{field:"taskNum",label:M("task.taskCount"),minWidth:15,formatter:(e,t,a)=>s(x,{type:"info"},(()=>a))},{field:"progress",label:M("task.taskProgress"),minWidth:40,formatter:(e,t,a)=>s(j,{percentage:a,type:"line",striped:!0,status:a<100?"":"success",stripedFlow:a<100})},{field:"creatTime",label:M("task.createTime"),minWidth:30},{field:"endTime",label:M("task.endTime"),minWidth:30,formatter:(e,t,a)=>""==a?"-":a},{field:"action",label:M("tableDemo.action"),minWidth:50,formatter:(e,t,a)=>{let o,r,u;const d=s(l,{type:"warning",onClick:()=>ye(e)},M("task.retest"));return i(n,null,[i(l,{type:"success",onClick:()=>de(e)},F(o=M("common.view"))?o:{default:()=>[o]}),d,i(l,{type:"danger",onClick:()=>pe(e)},F(r=M("common.delete"))?r:{default:()=>[r]}),i(l,{type:"primary",onClick:()=>G(e.id)},F(u=M("task.taskProgress"))?u:{default:()=>[u]})])}}]),J=a(!1);let O="";const G=async e=>{O=e,J.value=!0},K=()=>{J.value=!1},{tableRegister:q,tableState:B,tableMethods:Q}=C({fetchDataApi:async()=>{const{currentPage:e,pageSize:t}=B,a=await E(W.value,e.value,t.value);return{list:a.data.list,total:a.data.total}},immediate:!1}),{loading:X,dataList:Y,total:Z,currentPage:$,pageSize:ee}=B;ee.value=20;const{getList:te,getElTableExpose:ae}=Q;function oe(){return{background:"var(--el-fill-color-light)"}}const se=a(!1),le=async()=>{ie=M("task.addTask"),ue.value=!0,re.name="",re.target="",re.node=[],re.subdomainScan=!0,re.duplicates="None",re.subdomainConfig=[],re.urlScan=!0,re.sensitiveInfoScan=!0,re.pageMonitoring="JS",re.crawlerScan=!0,re.vulScan=!1,re.vulList=[],re.portScan=!1,re.ports="",re.dirScan=!0,re.waybackurl=!0,re.scheduledTasks=!0,re.hour=24,re.allNode=!1,se.value=!0};let ie=M("task.addTask");const ne=()=>{se.value=!1};let re=o({name:"",target:"",node:[],subdomainScan:!0,duplicates:"None",subdomainConfig:[],urlScan:!0,sensitiveInfoScan:!0,pageMonitoring:"JS",crawlerScan:!0,vulScan:!1,vulList:[],portScan:!1,ports:"",dirScan:!0,waybackurl:!0,scheduledTasks:!0,hour:24,allNode:!1}),ue=a(!0);const de=async e=>{const t=await P(e.id);if(200===t.code){const e=t.data;re.name=e.name,re.target=e.target,re.node=e.node,re.subdomainScan=e.subdomainScan,re.subdomainConfig=e.subdomainConfig,re.urlScan=e.urlScan,re.sensitiveInfoScan=e.sensitiveInfoScan,re.pageMonitoring=e.pageMonitoring,re.crawlerScan=e.crawlerScan,re.vulScan=e.vulScan,re.vulList=e.vulList,re.portScan=e.portScan,re.ports=e.ports,re.dirScan=e.dirScan,re.waybackurl=e.waybackurl,re.scheduledTasks=e.scheduledTasks,re.hour=e.hour,re.allNode=e.allNode,re.duplicates=e.duplicates}se.value=!0,ue.value=!1,ie=M("common.view")},ce=async()=>{window.confirm("Are you sure you want to delete the selected data?")&&await ve()},pe=async e=>{window.confirm("Are you sure you want to delete the selected data?")&&await ge(e)},me=a(!1),ge=async e=>{me.value=!0;try{await N([e.id]);me.value=!1,te()}catch(t){me.value=!1,te()}},fe=a([]),ve=async()=>{const e=await ae(),t=(null==e?void 0:e.getSelectionRows())||[];fe.value=t.map((e=>e.id)),me.value=!0;try{await N(fe.value);me.value=!1,te()}catch(a){me.value=!1,te()}},ye=async e=>{window.confirm("Are you sure you want to retest?")&&await he(e)},he=async e=>{try{await A(e.id),te()}catch(t){te()}};r((()=>{_e(),window.addEventListener("resize",_e)}));const ke=a(0),_e=()=>{const e=window.innerHeight||document.documentElement.clientHeight;ke.value=.75*e};return(t,a)=>(u(),d(n,null,[i(p(e),null,{default:c((()=>[i(p(b),null,{default:c((()=>[i(p(S),{span:1},{default:c((()=>[i(p(w),{class:"mx-1",style:{position:"relative",top:"8px"}},{default:c((()=>[m(g(p(M)("task.taskName"))+":",1)])),_:1})])),_:1}),i(p(S),{span:5},{default:c((()=>[i(p(f),{modelValue:W.value,"onUpdate:modelValue":a[0]||(a[0]=e=>W.value=e),placeholder:p(M)("common.inputText"),style:{height:"38px"}},null,8,["modelValue","placeholder"])])),_:1}),i(p(S),{span:5,style:{position:"relative",left:"16px"}},{default:c((()=>[i(p(v),{type:"primary",icon:p(k),style:{height:"100%"},onClick:H},{default:c((()=>[m("Search")])),_:1},8,["icon"])])),_:1})])),_:1}),i(p(b),null,{default:c((()=>[i(p(S),{style:{position:"relative",top:"16px"}},{default:c((()=>[y("div",V,[i(p(l),{type:"primary",onClick:le},{default:c((()=>[m(g(p(M)("task.addTask")),1)])),_:1}),i(p(l),{type:"danger",loading:me.value,onClick:ce},{default:c((()=>[m(g(p(M)("task.delTask")),1)])),_:1},8,["loading"])])])),_:1})])),_:1}),y("div",D,[i(p(T),{"tooltip-options":{offset:1,showArrow:!1,effect:"dark",enterable:!1,showAfter:0,popperOptions:{},popperClass:"test",placement:"bottom",hideAfter:0,disabled:!0},pageSize:p(ee),"onUpdate:pageSize":a[1]||(a[1]=e=>h(ee)?ee.value=e:null),currentPage:p($),"onUpdate:currentPage":a[2]||(a[2]=e=>h($)?$.value=e:null),columns:R,data:p(Y),stripe:"",border:!0,loading:p(X),"max-height":ke.value,resizable:!0,pagination:{total:p(Z),pageSizes:[20,30,50,100,200,500,1e3]},onRegister:p(q),headerCellStyle:oe,style:{fontFamily:"-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji"}},null,8,["pageSize","currentPage","columns","data","loading","max-height","pagination","onRegister"])])])),_:1}),i(p(L),{modelValue:se.value,"onUpdate:modelValue":a[3]||(a[3]=e=>se.value=e),title:p(ie),center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"}},{default:c((()=>[i(z,{closeDialog:ne,getList:p(te),vTaskForm:p(re),create:p(ue),taskid:"",schedule:!1},null,8,["getList","vTaskForm","create"])])),_:1},8,["modelValue","title"]),i(p(L),{modelValue:J.value,"onUpdate:modelValue":a[4]||(a[4]=e=>J.value=e),title:p(M)("task.taskProgress"),center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"},width:"70%","max-height":"700"},{default:c((()=>[i(U,{closeDialog:K,getProgressInfoID:p(O),getProgressInfotype:"scan",getProgressInforunnerid:""},null,8,["getProgressInfoID"])])),_:1},8,["modelValue","title"])],64))}});export{M as default};