import{_ as e}from"./ContentWrap.vue_vue_type_script_setup_true_lang-Cl5a1g_M.js";import{d as a,N as t,e as l,a1 as o,I as r,r as s,X as i,Q as n,E as p,Y as u,o as m,c as d,w as c,a as h,$ as f,f as g,t as x,a4 as _,H as b,a0 as v,l as y,_ as j}from"./index-Cr6AeRWq.js";import{_ as w}from"./Search.vue_vue_type_script_setup_true_lang-CBi3Eg7V.js";import{u as H}from"./useSearch-Fsyfplcb.js";import{u as k}from"./useTable-DixQRLyI.js";import{E as S}from"./el-card-f04OrW9f.js";import{E as z}from"./el-pagination-14Id82WI.js";import"./el-tag-CxICoOpO.js";import"./el-popper-BgAeKG77.js";import{E,a as P}from"./el-col-B0jdjJOD.js";import{_ as C,E as I,a as R}from"./Table.vue_vue_type_script_lang-BFDltahK.js";import"./el-checkbox-zfukSsNU.js";import"./el-tooltip-l0sNRNKZ.js";import{E as U}from"./el-text-C5s84vQY.js";import{E as V}from"./el-divider-gyCaIvGK.js";import{E as A}from"./el-link-Rl-0_Y9K.js";import{_ as W}from"./Dialog.vue_vue_type_style_index_0_lang-DNSEKIW8.js";import{u as M}from"./useCrudSchemas-CNKWpknJ.js";import{e as B}from"./index-DRni1RON.js";import"./useForm-5NmR2Rgd.js";import"./el-form-item-DAdIpZ2w.js";import"./castArray-DolOOpA2.js";import"./el-radio-group-DlaD2ltx.js";/* empty css                          */import"./el-input-number-DjDJ99ro.js";import"./el-select-v2-C2BLWs-e.js";import"./raf-FGCLxHfU.js";import"./useInput-8bUrr7Qz.js";import"./debounce-D2fEcvzH.js";import"./el-switch-Dv-A2hxv.js";import"./el-progress-DbNHGNbp.js";import"./InputPassword-CmeI5942.js";import"./style.css_vue_type_style_index_0_src_true_lang-819h2M_i.js";import"./JsonEditor.vue_vue_type_script_setup_true_lang-DUWpnlaB.js";import"./IconPicker-2jg_b_V6.js";/* empty css                   */import"./el-tab-pane-Bekx6uf7.js";import"./tsxHelper-B2K87hSo.js";import"./index-RQ44C1fT.js";import"./useIcon-CzeVxZTS.js";import"./el-image-viewer-CKf6yRwH.js";import"./el-dropdown-item-CjlYFbhA.js";import"./refs-COABADRW.js";import"./index-ukXxIniG.js";import"./tree-BfZhwLPs.js";import"./index-XBrIPdP2.js";const D={style:{color:"red"}},L={style:{whiteSpace:"pre-line"}};function N(e){return"function"==typeof e||"[object Object]"===Object.prototype.toString.call(e)&&!v(e)}const O=j(a({__name:"PageMonitoring",setup(a){const{t:v}=y(),{searchRegister:j}=H(),O=t([{field:"search",label:v("form.input"),component:"Input",formItemProps:{size:"large",style:{width:"100%"}},componentProps:{clearable:!1,slots:{suffix:()=>l(r,{class:"icon-button",onClick:Q,text:!0,style:"outline: none;background-color: transparent !important; color: inherit !important; box-shadow: none !important;position: relative;left: 24%"},{default:()=>[l(o,{icon:"tdesign:chat-bubble-help"},null)]})}}}]),q=[{operator:"=",meaning:v("searchHelp.like")},{operator:"!=",meaning:v("searchHelp.notIn")},{operator:"==",meaning:v("searchHelp.equal")},{operator:"&&",meaning:v("searchHelp.and")},{operator:"||",meaning:v("searchHelp.or")},{operator:"()",meaning:v("searchHelp.brackets")}],F=[{keyword:"url",example:'url="http://example.com"',explain:v("searchHelp.url")},{keyword:"hash",example:'hash="234658675623543"',explain:v("searchHelp.hash")},{keyword:"matched",example:'matched="https://example.com"',explain:v("searchHelp.matched")},{keyword:"diff",example:'diff="example"',explain:v("searchHelp.diff")},{keyword:"response",example:'response="root"',explain:v("searchHelp.response")},{keyword:"project",example:'project="Hackerone"',explain:v("searchHelp.project")}],T=s(!1),Q=()=>{T.value=!0},X=s(!0),Y=s("inline"),G=s("left"),J=s(""),K=e=>{J.value=e.search,ge()},Z=s(!1),$=t([{field:"index",label:v("tableDemo.index"),type:"index",minWidth:10},{field:"url",label:"url",minWidth:30,formatter:(e,a,t)=>l(A,{href:t,underline:!1},N(t)?t:{default:()=>[t]})},{field:"response1",label:v("PageMonitoring.oldResponseBody"),minWidth:30,formatter:(e,a,t)=>{let o;return l(n,null,[l(i,{type:"success",onClick:()=>le(e.response1,e.hash1)},N(o=v("common.view"))?o:{default:()=>[o]})])}},{field:"respone2",label:v("PageMonitoring.currentResponseBody"),minWidth:30,formatter:(e,a,t)=>{let o;return l(n,null,[l(i,{type:"success",onClick:()=>le(e.response2,e.hash2)},N(o=v("common.view"))?o:{default:()=>[o]})])}},{field:"diff",label:"diff",formatter:(e,a,t)=>{const o=t.split("\n").map(((e,a)=>l("div",{key:a},[e])));return l(p,{minSize:10,maxHeight:200},{default:()=>[l("div",{class:"scrollbar-demo-item"},[o])]})}},{field:"action",label:v("tableDemo.action"),minWidth:30,formatter:(e,a,t)=>{let o;return l(n,null,[l(i,{type:"success",onClick:()=>se(e.history_diff)},N(o=v("asset.historyDiff"))?o:{default:()=>[o]})])}}]),ee=s(!1),ae=s(""),te=s(""),le=(e,a)=>{ee.value=!0,ae.value=e,te.value=a},oe=s(!1),re=s([]),se=e=>{re.value=e,oe.value=!0},{allSchemas:ie}=M($),{tableRegister:ne,tableState:pe,tableMethods:ue}=k({fetchDataApi:async()=>{const{currentPage:e,pageSize:a}=pe,t=await B(J.value,e.value,a.value);return{list:t.data.list,total:t.data.total}}}),{loading:me,dataList:de,total:ce,currentPage:he,pageSize:fe}=pe,{getList:ge}=ue;function xe(){return{background:"var(--el-fill-color-light)"}}u((()=>{be(),window.addEventListener("resize",be)}));const _e=s(0),be=()=>{const e=window.innerHeight||document.documentElement.clientHeight;_e.value=.7*e};return(a,t)=>(m(),d(n,null,[l(h(e),{style:{height:"80px"}},{default:c((()=>[l(h(w),{schema:O,"is-col":X.value,layout:Y.value,"show-reset":!1,"button-position":G.value,"search-loading":Z.value,onSearch:K,onReset:K,onRegister:h(j)},null,8,["schema","is-col","layout","button-position","search-loading","onRegister"])])),_:1}),l(h(P),null,{default:c((()=>[l(h(E),null,{default:c((()=>[l(h(S),null,{default:c((()=>[l(h(C),{pageSize:h(fe),"onUpdate:pageSize":t[0]||(t[0]=e=>f(fe)?fe.value=e:null),currentPage:h(he),"onUpdate:currentPage":t[1]||(t[1]=e=>f(he)?he.value=e:null),columns:h(ie).tableColumns,data:h(de),"max-height":_e.value,stripe:"",border:!0,loading:h(me),resizable:!0,onRegister:h(ne),headerCellStyle:xe,tooltipOptions:{disabled:!0,showArrow:!1,effect:"dark",enterable:!1,offset:0,placement:"top",popperClass:"",popperOptions:{},showAfter:0,hideAfter:0},style:{fontFamily:"-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji"}},null,8,["pageSize","currentPage","columns","data","max-height","loading","onRegister"])])),_:1})])),_:1}),l(h(E),{":span":24},{default:c((()=>[l(h(S),null,{default:c((()=>[l(h(z),{pageSize:h(fe),"onUpdate:pageSize":t[2]||(t[2]=e=>f(fe)?fe.value=e:null),currentPage:h(he),"onUpdate:currentPage":t[3]||(t[3]=e=>f(he)?he.value=e:null),"page-sizes":[10,20,50,100,200,500,1e3],layout:"total, sizes, prev, pager, next, jumper",total:h(ce)},null,8,["pageSize","currentPage","total"])])),_:1})])),_:1})])),_:1}),l(h(W),{modelValue:ee.value,"onUpdate:modelValue":t[4]||(t[4]=e=>ee.value=e),title:"ResponseBody",center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"},width:"70%","max-height":_e.value},{default:c((()=>[l(h(p),{"max-height":_e.value},{default:c((()=>[g("div",D,"Hash: "+x(te.value),1),l(h(V)),g("div",L,x(ae.value),1)])),_:1},8,["max-height"])])),_:1},8,["modelValue","max-height"]),l(h(W),{modelValue:oe.value,"onUpdate:modelValue":t[5]||(t[5]=e=>oe.value=e),title:"Historical changes",center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"},width:"70%","max-height":_e.value},{default:c((()=>[g("div",null,[(m(!0),d(n,null,_(re.value,((e,a)=>(m(),d("div",{key:a,style:{whiteSpace:"pre-line"}},[l(h(U),null,{default:c((()=>[b(x(e),1)])),_:2},1024),l(h(V),{style:{background:"#e99696"}})])))),128))])])),_:1},8,["modelValue","max-height"]),l(h(W),{modelValue:T.value,"onUpdate:modelValue":t[6]||(t[6]=e=>T.value=e),title:h(v)("common.querysyntax"),center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"}},{default:c((()=>[l(h(P),null,{default:c((()=>[l(h(E),null,{default:c((()=>[l(h(U),{tag:"b",size:"small"},{default:c((()=>[b(x(h(v)("searchHelp.operator")),1)])),_:1}),l(h(V),{direction:"vertical"}),l(h(U),{size:"small",type:"danger"},{default:c((()=>[b(x(h(v)("searchHelp.notice")),1)])),_:1})])),_:1}),l(h(E),{style:{"margin-top":"10px"}},{default:c((()=>[l(h(I),{headerCellStyle:xe,data:q},{default:c((()=>[l(h(R),{prop:"operator",label:h(v)("searchHelp.operator"),width:"300"},null,8,["label"]),l(h(R),{prop:"meaning",label:h(v)("searchHelp.meaning")},null,8,["label"])])),_:1})])),_:1}),l(h(E),{style:{"margin-top":"15px"}},{default:c((()=>[l(h(U),{tag:"b",size:"small"},{default:c((()=>[b(x(h(v)("searchHelp.keywords")),1)])),_:1})])),_:1}),l(h(E),{style:{"margin-top":"10px"}},{default:c((()=>[l(h(I),{headerCellStyle:xe,data:F},{default:c((()=>[l(h(R),{prop:"keyword",label:h(v)("searchHelp.keywords")},null,8,["label"]),l(h(R),{prop:"example",label:h(v)("searchHelp.example")},null,8,["label"]),l(h(R),{prop:"explain",label:h(v)("searchHelp.explain")},null,8,["label"])])),_:1})])),_:1})])),_:1})])),_:1},8,["modelValue","title"])],64))}}),[["__scopeId","data-v-d344915e"]]);export{O as default};