import{Y as e,d as a,bg as t,bL as s,a7 as l,r as o,a6 as r,aA as d,b6 as i,o as n,i as f,e as c,w as u,Q as p,a as v,dr as b,c0 as y,f as h,aJ as m,ae as k,n as A,c as _,a8 as C,t as w,j as g,A as L,bd as F,af as R,h as E,c1 as x,a9 as $,ag as S}from"./index-ChGT_YCB.js";import{d as z,a as B,u as T}from"./Dialog.vue_vue_type_style_index_0_lang-CjT1Och0.js";const D=e({...z,direction:{type:String,default:"rtl",values:["ltr","rtl","ttb","btt"]},size:{type:[String,Number],default:"30%"},withHeader:{type:Boolean,default:!0},modalFade:{type:Boolean,default:!0},headerAriaLevel:{type:String,default:"2"}}),I=B,P=["aria-label","aria-labelledby","aria-describedby"],j=["id","aria-level"],q=["aria-label"],H=["id"],O=a({name:"ElDrawer",inheritAttrs:!1});const J=S($(a({...O,props:D,emits:I,setup(e,{expose:a}){const $=e,S=t();s({scope:"el-drawer",from:"the title slot",replacement:"the header slot",version:"3.0.0",ref:"https://element-plus.org/en-US/component/drawer.html#slots"},l((()=>!!S.title)));const z=o(),B=o(),D=r("drawer"),{t:I}=d(),{afterEnter:O,afterLeave:J,beforeLeave:M,visible:N,rendered:Q,titleId:U,bodyId:Y,zIndex:G,onModalClick:K,onOpenAutoFocus:V,onCloseAutoFocus:W,onFocusoutPrevented:X,onCloseRequested:Z,handleClose:ee}=T($,z),ae=l((()=>"rtl"===$.direction||"ltr"===$.direction)),te=l((()=>i($.size)));return a({handleClose:ee,afterEnter:O,afterLeave:J}),(e,a)=>(n(),f(x,{to:"body",disabled:!e.appendToBody},[c(E,{name:v(D).b("fade"),onAfterEnter:v(O),onAfterLeave:v(J),onBeforeLeave:v(M),persisted:""},{default:u((()=>[p(c(v(b),{mask:e.modal,"overlay-class":e.modalClass,"z-index":v(G),onClick:v(K)},{default:u((()=>[c(v(y),{loop:"",trapped:v(N),"focus-trap-el":z.value,"focus-start-el":B.value,onFocusAfterTrapped:v(V),onFocusAfterReleased:v(W),onFocusoutPrevented:v(X),onReleaseRequested:v(Z)},{default:u((()=>[h("div",m({ref_key:"drawerRef",ref:z,"aria-modal":"true","aria-label":e.title||void 0,"aria-labelledby":e.title?void 0:v(U),"aria-describedby":v(Y)},e.$attrs,{class:[v(D).b(),e.direction,v(N)&&"open"],style:v(ae)?"width: "+v(te):"height: "+v(te),role:"dialog",onClick:a[1]||(a[1]=k((()=>{}),["stop"]))}),[h("span",{ref_key:"focusStartRef",ref:B,class:A(v(D).e("sr-focus")),tabindex:"-1"},null,2),e.withHeader?(n(),_("header",{key:0,class:A(v(D).e("header"))},[e.$slots.title?C(e.$slots,"title",{key:1},(()=>[g(" DEPRECATED SLOT ")])):C(e.$slots,"header",{key:0,close:v(ee),titleId:v(U),titleClass:v(D).e("title")},(()=>[e.$slots.title?g("v-if",!0):(n(),_("span",{key:0,id:v(U),role:"heading","aria-level":e.headerAriaLevel,class:A(v(D).e("title"))},w(e.title),11,j))])),e.showClose?(n(),_("button",{key:2,"aria-label":v(I)("el.drawer.close"),class:A(v(D).e("close-btn")),type:"button",onClick:a[0]||(a[0]=(...e)=>v(ee)&&v(ee)(...e))},[c(v(L),{class:A(v(D).e("close"))},{default:u((()=>[c(v(F))])),_:1},8,["class"])],10,q)):g("v-if",!0)],2)):g("v-if",!0),v(Q)?(n(),_("div",{key:1,id:v(Y),class:A(v(D).e("body"))},[C(e.$slots,"default")],10,H)):g("v-if",!0),e.$slots.footer?(n(),_("div",{key:2,class:A(v(D).e("footer"))},[C(e.$slots,"footer")],2)):g("v-if",!0)],16,P)])),_:3},8,["trapped","focus-trap-el","focus-start-el","onFocusAfterTrapped","onFocusAfterReleased","onFocusoutPrevented","onReleaseRequested"])])),_:3},8,["mask","overlay-class","z-index","onClick"]),[[R,v(N)]])])),_:3},8,["name","onAfterEnter","onAfterLeave","onBeforeLeave"])],8,["disabled"]))}}),[["__file","drawer.vue"]]));export{J as E};