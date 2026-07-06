(function(){let e=document.createElement(`link`).relList;if(e&&e.supports&&e.supports(`modulepreload`))return;for(let e of document.querySelectorAll(`link[rel="modulepreload"]`))n(e);new MutationObserver(e=>{for(let t of e)if(t.type===`childList`)for(let e of t.addedNodes)e.tagName===`LINK`&&e.rel===`modulepreload`&&n(e)}).observe(document,{childList:!0,subtree:!0});function t(e){let t={};return e.integrity&&(t.integrity=e.integrity),e.referrerPolicy&&(t.referrerPolicy=e.referrerPolicy),e.crossOrigin===`use-credentials`?t.credentials=`include`:e.crossOrigin===`anonymous`?t.credentials=`omit`:t.credentials=`same-origin`,t}function n(e){if(e.ep)return;e.ep=!0;let n=t(e);fetch(e.href,n)}})();var e=1,t=2,n=4,r=8,i=16;e|t|n|r|i|32;function a(e,t){return e|t}var o=r;function s(){let e=o;return o=e===r?i:r,e}var c=n,l=!1,u=[];function d(){return c}function f(e){let t=c,n=l;c=s(),l=!0;try{e()}finally{c=t,l=n,ee()}}function ee(){let e=u;u=[];for(let t of e)t()}var te=new Map,ne=0,re=!0;function ie(e){if(!re)return``;let t=te.get(e);t||(t={name:e,index:ne++,instanceCount:0},te.set(e,t),ae());let n=t.instanceCount++;return`s-${t.index}-${n}`}function ae(){if(typeof globalThis<`u`){let e={};for(let t of te.values())e[t.index]=t.name;globalThis.__SPECIFYJS_COMPONENTS__=e}}function p(e){return typeof e==`string`?e:typeof e==`function`?e.name||`Anonymous`:e==null?`Unknown`:`Component`}function m(e){e()}function oe(e){typeof queueMicrotask==`function`?queueMicrotask(e):Promise.resolve().then(e)}var se=Symbol.for(`spec.element`),ce=Symbol.for(`spec.fragment`),le=Symbol.for(`spec.portal`),ue=Symbol.for(`spec.provider`),de=Symbol.for(`spec.consumer`),fe=Symbol.for(`spec.forward_ref`),pe=Symbol.for(`spec.memo`),me=Symbol.for(`spec.suspense`),h=Symbol.for(`spec.strict_mode`),he=Symbol.for(`spec.profiler`);function g(e){return typeof e==`object`&&!!e&&e.$$typeof===se}function _(e,t,...n){let r=null,i=null,a={};if(t!=null){t.key!==void 0&&(r=``+t.key),t.ref!==void 0&&(i=t.ref);for(let e in t)Object.prototype.hasOwnProperty.call(t,e)&&e!==`key`&&e!==`ref`&&e!==`__proto__`&&e!==`constructor`&&e!==`prototype`&&(a[e]=t[e])}if(n.length===1?a.children=n[0]:n.length>1&&(a.children=n),typeof e==`function`&&e.defaultProps){let t=e.defaultProps;for(let e of Object.keys(t))e!==`__proto__`&&e!==`constructor`&&a[e]===void 0&&(a[e]=t[e])}return{$$typeof:se,type:e,props:a,key:r,ref:i}}var ge=null;function _e(e){ge=e}function v(){if(ge===null)throw Error(`Invalid hook call. Hooks can only be called inside the body of a function component.`);return ge}function ve(e){return v().useState(e)}function ye(e,t){return v().useEffect(e,t)}function be(e){return v().useContext(e)}function xe(e,t){return v().useCallback(e,t)}function Se(e,t){return v().useMemo(e,t)}function Ce(e){return v().useRef(e)}function y(e){let t=Ce(e);t.current=e,ye(()=>{let e=t.current,n=[],r=document.title;if(e.title&&(document.title=e.title,n.push(()=>{document.title=r})),e.description&&n.push(b(`name`,`description`,e.description)),e.keywords&&n.push(b(`name`,`keywords`,e.keywords)),e.author&&n.push(b(`name`,`author`,e.author)),e.canonical){let t=document.createElement(`link`);t.rel=`canonical`,t.href=e.canonical,document.head.appendChild(t),n.push(()=>t.remove())}if(e.og)for(let[t,r]of Object.entries(e.og))n.push(b(`property`,`og:${t}`,r));if(e.twitter)for(let[t,r]of Object.entries(e.twitter))n.push(b(`name`,`twitter:${t}`,r));if(e.httpEquiv){let t={csp:`Content-Security-Policy`,referrer:`Referrer-Policy`,contentType:`X-Content-Type-Options`,frameOptions:`X-Frame-Options`,cacheControl:`Cache-Control`};for(let[r,i]of Object.entries(e.httpEquiv)){if(i===void 0)continue;let e=t[r]??r;n.push(Te(e,i))}}if(e.meta)for(let t of e.meta)t.name?n.push(b(`name`,t.name,t.content)):t.property&&n.push(b(`property`,t.property,t.content));return()=>{for(let e of n)e()}},[JSON.stringify([e.title,e.description,e.keywords,e.author,e.canonical,e.og,e.twitter,e.httpEquiv,e.meta])])}function we(e){return e.replace(/["\\]/g,`\\$&`)}function b(e,t,n){let r=`meta[${e}="${we(t)}"]`,i=document.querySelector(r),a=i!==null,o=i?.content;return i||(i=document.createElement(`meta`),i.setAttribute(e,t),document.head.appendChild(i)),i.content=n,()=>{a&&o!==void 0?i.content=o:a||i.remove()}}function Te(e,t){let n=`meta[http-equiv="${we(e)}"]`,r=document.querySelector(n),i=r!==null,a=r?.content;return r||(r=document.createElement(`meta`),r.setAttribute(`http-equiv`,e),document.head.appendChild(r)),r.content=t,()=>{i&&a!==void 0?r.content=a:i||r.remove()}}function Ee(e){let t={$$typeof:Symbol.for(`spec.context`),Provider:null,Consumer:null,_currentValue:e,_defaultValue:e},n={$$typeof:ue,_context:t},r={$$typeof:de,_context:t};return t.Provider=n,t.Consumer=r,t}var De=class{constructor(e){this._fiber=null,this._pendingState=[],this._forceUpdate=!1,this.props=e,this.state={},this.context=void 0}setState(e,t){this._pendingState.push(e),this._enqueueUpdate(t)}forceUpdate(e){this._forceUpdate=!0,this._enqueueUpdate(e)}render(){return null}_enqueueUpdate(e){}};De.prototype.isSpecComponent=!0;var Oe=class extends De{shouldComponentUpdate(e,t){return!ke(this.props,e)||!ke(this.state,t)}};Oe.prototype.isPureSpecComponent=!0;function ke(e,t){if(Object.is(e,t))return!0;if(typeof e!=`object`||!e||typeof t!=`object`||!t)return!1;let n=Object.keys(e),r=Object.keys(t);if(n.length!==r.length)return!1;let i=e,a=t;for(let e of n)if(!Object.prototype.hasOwnProperty.call(a,e)||!Object.is(i[e],a[e]))return!1;return!0}var x=Ee({pathname:`/`,params:{},navigate:()=>{throw Error(`useNavigate must be used inside a <Router> component.`)},basePath:``});function Ae(){let e=(typeof window<`u`?window.location.hash:``).replace(/^#\/?/,`/`);return e===``?`/`:e}var je={pathname:Ae(),hash:typeof window<`u`?window.location.hash:``},S=new Set;function Me(){je={pathname:Ae(),hash:typeof window<`u`?window.location.hash:``};for(let e of S)e()}typeof window<`u`&&window.addEventListener(`hashchange`,Me);function Ne(e){return S.add(e),()=>{S.delete(e)}}function C(){return je}function Pe(e,t){let n=e.startsWith(`#`)?e:`#`+e;if(t?.replace){let e=window.location.pathname+window.location.search+n;window.history.replaceState(null,``,e)}else window.location.hash=n;Me()}function Fe(e){let[t,n]=ve(()=>C().pathname);ye(()=>(n(C().pathname),Ne(()=>{n(C().pathname)})),[]);let r=xe((...e)=>{Pe(e[0],e[1])},[]),i=Se(()=>({pathname:t,params:{},navigate:r,basePath:``}),[t,r]);return _(x.Provider,{value:i},e.children)}function Ie(e){let t=e.length;for(;t>1&&e.charCodeAt(t-1)===47;)t--;return t===e.length?e:e.slice(0,t)}function Le(e,t,n){let r=n?.exact??!1,i=e===`/`?`/`:Ie(e),a=t===`/`?`/`:Ie(t);if(i===`/`){let t=a===`/`;return r&&!t?null:{params:{},isExact:t,path:e,url:`/`}}let o=i.split(`/`).filter(Boolean),s=a.split(`/`).filter(Boolean),c=o[o.length-1]===`*`;if(!c&&r&&s.length!==o.length||!c&&s.length<o.length)return null;let l={},u=[];for(let t=0;t<o.length;t++){let n=o[t];if(n===`*`)return l[`*`]=s.slice(t).join(`/`),{params:l,isExact:!0,path:e,url:`/`+s.join(`/`)};if(t>=s.length)return null;let r=s[t];if(n.startsWith(`:`)){let e=n.slice(1);l[e]=decodeURIComponent(r),u.push(r);continue}if(n.toLowerCase()!==r.toLowerCase())return null;u.push(r)}let d=`/`+u.join(`/`),f=s.length===o.length;return r&&!f?null:{params:l,isExact:f,path:e,url:d}}function w(e){let t=be(x),n=Le(t.basePath+e.path,t.pathname,{exact:e.exact??!1}),r=Se(()=>({pathname:t.pathname,params:{...t.params,...n?n.params:{}},navigate:t.navigate,basePath:n?n.url:t.basePath}),[t.pathname,t.navigate,t.params,t.basePath,n]);if(!n)return null;for(let e of Object.keys(t.params))e in r.params||(r.params[e]=t.params[e]);let i=e.component?_(e.component,{...r.params}):e.children;return _(x.Provider,{value:r},i)}function T(e){let{to:t,className:n,activeClassName:r,exact:i,children:a,...o}=e,s=be(x),c=Le(t,s.pathname,{exact:i??!1})!==null,l=xe(((...e)=>{e[0].preventDefault(),s.navigate(t)}),[t,s.navigate]),u=[n,c?r:null].filter(Boolean).join(` `)||void 0;return _(`a`,{...o,href:`#`+t,onClick:l,className:u},a)}function Re(){return be(x)}Ee({flags:{},isEnabled:()=>!1,setFlag:()=>{},loading:!1});function ze(e){if(typeof e==`string`)return 3;if(typeof e==`function`)return+!!e.prototype?.isSpecComponent;if(typeof e==`symbol`){if(e===ce)return 5;if(e===me)return 11;if(e===h)return 5;if(e===he)return 12;if(e===le)return 13}if(typeof e==`object`&&e){let t=e.$$typeof;if(t===fe)return 8;if(t===pe)return 9;if(t===ue)return 6;if(t===de)return 7}return 3}function E(e,t=0){return O(ze(e.type),e.type,e.key,e.ref,e.props,t)}function D(e,t=0){return O(4,null,null,null,{text:e},t)}function Be(){return O(2,null,null,null,{},0)}function O(e,t,n,r,i,a){return{tag:e,type:t,key:n,ref:r,stateNode:null,pendingProps:i,memoizedProps:null,memoizedState:null,return:null,child:null,sibling:null,index:0,alternate:null,effectTag:0,updateQueue:null,dependencies:null,lanes:a,childLanes:0}}function Ve(e,t){let n=e.alternate;return n===null?(n=O(e.tag,e.type,e.key,e.ref,t,e.lanes),n.stateNode=e.stateNode,n.alternate=e,e.alternate=n):(n.pendingProps=t,n.effectTag=0,n.child=null,n.sibling=null),n.memoizedProps=e.memoizedProps,n.memoizedState=e.memoizedState,n.updateQueue=e.updateQueue,n.lanes=e.lanes,n.childLanes=e.childLanes,n}function He(e){let t=[],n=[e];for(;n.length>0;){let e=n.pop();if(!(e==null||typeof e==`boolean`)){if(Array.isArray(e)){for(let t=e.length-1;t>=0;t--)n.push(e[t]);continue}(g(e)||typeof e==`string`||typeof e==`number`)&&t.push(e)}}return t}function k(e,t,n,r){let i=He(n);return i.length===0?(N(e,t),null):i.length===1&&!Array.isArray(i[0])?Ue(e,t,i[0],r):Ge(e,t,i,r)}function Ue(e,t,n,r){if(typeof n==`string`||typeof n==`number`)return We(e,t,n,r);let i=n.key,a=t;for(;a!==null;){if(a.key===i){if(Ke(a,n)){N(e,a.sibling);let t=A(a,n.props);return t.ref=n.ref,t.return=e,t}N(e,a);break}else M(e,a);a=a.sibling}let o=E(n,r);return o.ref=n.ref,o.return=e,o.effectTag=1,o}function We(e,t,n,r){if(t!==null&&t.tag===4){N(e,t.sibling);let r=A(t,{text:n});return r.return=e,r}N(e,t);let i=D(n,r);return i.return=e,i.effectTag=1,i}function Ge(e,t,n,r){let i=null,a=null,o=t,s=0,c=0,l=null;for(;o!==null&&s<n.length;s++){o.index>s?(l=o,o=null):l=o.sibling;let t=n[s],u=Je(e,o,t,r);if(u===null){o===null&&(o=l);break}o&&u.alternate===null&&M(e,o),c=j(u,c,s),a===null?i=u:a.sibling=u,a=u,o=l}if(s===n.length)return N(e,o),i;if(o===null){for(;s<n.length;s++){let t=Ze(e,n[s],r);t!==null&&(c=j(t,c,s),a===null?i=t:a.sibling=t,a=t)}return i}let u=Qe(o);for(;s<n.length;s++){let t=n[s],o=$e(e,u,t,s,r);if(o!==null){if(o.alternate!==null){let e=o.key===null?s:o.key;u.delete(e)}c=j(o,c,s),a===null?i=o:a.sibling=o,a=o}}return u.forEach(t=>{M(e,t)}),i}function Ke(e,t){return e.type===t.type}function A(e,t){let n=qe(e,t);return n.index=0,n.sibling=null,n}function qe(e,t){let n=e.alternate;return n===null?(n={...e,pendingProps:t,effectTag:0,child:null,sibling:null,alternate:e},e.alternate=n):(n.pendingProps=t,n.effectTag=0,n.child=null,n.sibling=null),n.memoizedProps=e.memoizedProps,n.memoizedState=e.memoizedState,n.updateQueue=e.updateQueue,n}function Je(e,t,n,r){let i=t===null?null:t.key;if(typeof n==`string`||typeof n==`number`)return i===null?Ye(e,t,n,r):null;if(g(n)){let a=n;return a.key===i?Xe(e,t,a,r):null}return null}function Ye(e,t,n,r){if(t===null||t.tag!==4){let t=D(n,r);return t.return=e,t.effectTag=1,t}let i=A(t,{text:n});return i.return=e,i}function Xe(e,t,n,r){if(t!==null&&t.type===n.type){let r=A(t,n.props);return r.ref=n.ref,r.return=e,r}let i=E(n,r);return i.ref=n.ref,i.return=e,i.effectTag=1,i}function Ze(e,t,n){if(typeof t==`string`||typeof t==`number`){let r=D(t,n);return r.return=e,r.effectTag=1,r}if(g(t)){let r=E(t,n);return r.ref=t.ref,r.return=e,r.effectTag=1,r}return null}function j(e,t,n){e.index=n;let r=e.alternate;if(r!==null){let n=r.index;return n<t?(e.effectTag=1,t):n}return e.effectTag=1,t}function M(e,t){t.effectTag=4,e.updateQueue||=[],e.updateQueue.push(t)}function N(e,t){let n=t;for(;n!==null;)M(e,n),n=n.sibling}function Qe(e){let t=new Map,n=e;for(;n!==null;){let e=n.key===null?n.index:n.key;t.set(e,n),n=n.sibling}return t}function $e(e,t,n,r,i){if(typeof n==`string`||typeof n==`number`)return Ye(e,t.get(r)||null,n,i);if(g(n)){let a=n,o=a.key===null?r:a.key;return Xe(e,t.get(o)||null,a,i)}return null}var P=null,F=null,I=null,L=0,R=null,z=null;function B(e){P=e,F=null,I=null,L=0,R=null,z=null}function et(){return P}function tt(){return L}function V(){let e;if(P===null)throw Error(`Hooks can only be called inside a function component.`);let t=P.alternate;if(t!==null){if(I=I===null?t.memoizedState??null:I.next,I===null)throw Error(`Rendered more hooks than during the previous render.`);e={memoizedState:I.memoizedState,queue:I.queue,next:null}}else e={memoizedState:null,queue:null,next:null};return F===null?(P.memoizedState=e,F=e):(F.next=e,F=e),L++,e}function nt(e,t,n,r){let i={tag:e,create:t,destroy:n,deps:r,next:null};return z===null?(R=i,z=i):(z.next=i,z=i),i}function H(){return R}function rt(e,t){if(e===void 0||t===void 0||e.length!==t.length)return!1;for(let n=0;n<e.length;n++)if(!Object.is(e[n],t[n]))return!1;return!0}var it=1e3,at=new Set,ot=typeof process<`u`&&!1;function U(e){ot||at.has(e)||at.size>=it||(at.add(e),typeof console<`u`&&console.warn(`[SpecifyJS] ${e}`))}function st(e){ot||typeof console<`u`&&console.error(`[SpecifyJS] ${e}`)}var ct=50,lt=25,ut=5,dt=new WeakMap,ft=0;function pt(){ft++}function mt(e,t){let n=dt.get(e);if((!n||n.cycleId!==ft)&&(n={count:0,cycleId:ft},dt.set(e,n)),n.count++,n.count>ct){let e=`Maximum update depth exceeded. Component "${t}" re-rendered ${n.count} times in a single cycle. This usually means a useEffect or setState call is creating an infinite loop. Check that useEffect dependencies are stable (use useMemo/useCallback for objects/arrays).`;throw st(e),n.count=0,Error(`[SpecifyJS] ${e}`)}}var ht=new WeakMap,gt=0;function _t(){gt++}function vt(e,t){let n=ht.get(e);(!n||n.commitId!==gt)&&(n={count:0,commitId:gt},ht.set(e,n)),n.count++,n.count>lt&&U(`Effect cycle detected in "${t}". Effects fired ${n.count} times in a single commit. This may indicate an effect that triggers a state update on every run. Ensure effect dependencies are correct and stable.`)}var yt=new WeakMap;function bt(e,t,n,r,i){if(!n)return;let a=yt.get(e);a||(a=new Map,yt.set(e,a));let o=a.get(t);o||(o={consecutiveChanges:0,lastDepsKey:``,warned:!1},a.set(t,o)),r?o.consecutiveChanges++:o.consecutiveChanges=0,o.consecutiveChanges>=ut&&!o.warned&&(o.warned=!0,U(`Unstable useEffect dependencies in "${i}" (hook #${t}). Dependencies changed on ${o.consecutiveChanges} consecutive renders. This may cause unnecessary effect re-runs. Wrap objects/arrays in useMemo and functions in useCallback to stabilize references.`))}var xt=new WeakMap,W=0;function St(){W++}function Ct(e,t){let n=xt.get(e);(!n||n.frameId!==W)&&(n={count:0,frameId:W},xt.set(e,n)),n.count++,n.count===100&&U(`Rapid state updates in "${t}": setState called ${n.count} times in a single frame. This may indicate a loop or missing dependency optimization.`)}var G=null;function wt(e){G=e}function Tt(){let e=et();if(!e)throw Error(`Invalid hook call.`);return e}function Et(e,t){e.lanes=a(e.lanes,t);let n=e.return;for(;n!==null;)n.childLanes=a(n.childLanes,t),n=n.return}function Dt(e){let t=V();t.queue===null&&(t.memoizedState=typeof e==`function`?e():e,t.queue=[]);let n=t.queue;if(n.length>0){let e=t.memoizedState;for(let t of n){let n=t.action;e=typeof n==`function`?n(e):n}t.memoizedState=e,n.length=0}let r=t.memoizedState,i=Tt();return[r,e=>{let t=(typeof i.type==`function`?p(i.type)||i.type.name:``)||`Anonymous`;Ct(i,t);let r=d();n.length>=1e4&&(typeof console<`u`&&console.warn(`[SpecifyJS] Hook update queue exceeded 10000 — dropping oldest updates`),n.splice(0,n.length-5e3)),n.push({action:e}),Et(i,r),G&&m(()=>G(i))}]}function Ot(e,t,n){let r=V();r.queue===null&&(r.memoizedState=n?n(t):t,r.queue=[]);let i=r.queue;if(i.length>0){let t=r.memoizedState;for(let n of i)t=e(t,n.action);r.memoizedState=t,i.length=0}let a=r.memoizedState,o=Tt();return[a,e=>{let t=(typeof o.type==`function`?p(o.type)||o.type.name:``)||`Anonymous`;Ct(o,t);let n=d();i.push({action:e}),Et(o,n),G&&m(()=>G(o))}]}function K(e,t,n){let r=V(),i=tt(),a=r.memoizedState,o=a?.effect?.destroy??null;if(a!==null&&n!==void 0){let e=rt(a.deps,n),r=et();if(r){let t=(typeof r.type==`function`?p(r.type)||r.type.name:``)||`Anonymous`;bt(r,i,n,!e,t)}if(e){nt(0,t,o,n);return}}r.memoizedState={deps:n,effect:nt(1|e,t,o,n)}}function kt(e,t){K(4,e,t)}function At(e,t){K(2,e,t)}function jt(e,t){K(8,e,t)}function Mt(e){return V(),e._currentValue}function Nt(e,t){let n=V(),r=n.memoizedState;return r!==null&&rt(r[1],t)?r[0]:(n.memoizedState=[e,t],e)}function Pt(e,t){let n=V(),r=n.memoizedState;if(r!==null&&rt(r[1],t))return r[0];let i=e();return n.memoizedState=[i,t],i}function Ft(e){let t=V();return t.memoizedState===null&&(t.memoizedState={current:e}),t.memoizedState}function It(e,t,n){K(2,()=>{let n=t();if(e!==null)return typeof e==`function`?(e(n),()=>e(null)):(e.current=n,()=>{e.current=null})},n)}function Lt(e,t){}var Rt=0;function zt(){let e=V();return e.memoizedState===null&&(e.memoizedState=`:l${Rt++}:`),e.memoizedState}function Bt(e){let[t,n]=Dt(e);return K(4,()=>{f(()=>{n(e)})},[e]),t}function Vt(){let[e,t]=Dt(!1);return[e,Nt((e=>{t(!0),f(()=>{t(!1),e()})}),[t])]}function Ht(e,t,n){let r=V(),i=Tt(),a=t();return r.memoizedState=a,K(4,()=>e(()=>{let e=t();Object.is(r.memoizedState,e)||(r.memoizedState=e,G&&m(()=>G(i)))}),[e,t]),a}var Ut={useState:Dt,useEffect:kt,useContext:Mt,useReducer:Ot,useCallback:Nt,useMemo:Pt,useRef:Ft,useImperativeHandle:It,useLayoutEffect:At,useDebugValue:Lt,useId:zt,useDeferredValue:Bt,useTransition:Vt,useSyncExternalStore:Ht,useInsertionEffect:jt};function Wt(){_e(Ut)}function Gt(){_e(null)}var Kt=null;function qt(e){Kt?.onCommitFiberRoot?.(e)}function Jt(){return typeof performance<`u`&&typeof performance.now==`function`?performance.now():Date.now()}var Yt=5;function Xt(){Jt()+Yt}var q=null;function Zt(){if(q!==null){Xt();let e=q,t=!1;try{let n=e();typeof n==`function`?(q=n,t=!0):q=null}catch(e){throw q=null,e}t&&Qt()}}var Qt;if(typeof MessageChannel<`u`){let e=new MessageChannel;e.port1.onmessage=Zt,Qt=()=>{e.port2.postMessage(null)}}else Qt=()=>{setTimeout(Zt,0)};function $t(e,t){e===`button`&&!t.children&&!t[`aria-label`]&&!t[`aria-labelledby`]&&U(`<button> without text content or aria-label. Add text, aria-label, or aria-labelledby for accessibility.`),e===`img`&&t.alt===void 0&&U(`<img> without alt attribute. Add alt="" for decorative images or descriptive text for informative images.`),(e===`input`||e===`select`||e===`textarea`)&&!t[`aria-label`]&&!t[`aria-labelledby`]&&!t.id&&U(`<${e}> without aria-label, aria-labelledby, or id (for label[for]). Form controls need accessible names.`),(e===`div`||e===`span`)&&t.onClick&&!t.role&&U(`<${e}> with onClick but no role attribute. Add role="button" and tabIndex={0} for keyboard accessibility.`),t.role===`button`&&t.tabIndex===void 0&&e!==`button`&&U(`Element with role="button" but no tabIndex. Add tabIndex={0} to make it keyboard-focusable.`),t.role===`dialog`&&!t[`aria-label`]&&!t[`aria-labelledby`]&&U(`Element with role="dialog" needs aria-label or aria-labelledby.`),e===`a`&&t.target===`_blank`&&!t.rel&&U(`<a target="_blank"> without rel attribute. Add rel="noopener noreferrer" for security.`),t.role===`tablist`&&!t[`aria-label`]&&!t[`aria-labelledby`]&&U(`Element with role="tablist" needs aria-label or aria-labelledby.`)}var en=new Map;function tn(e){let t=Be();return t.stateNode=e,{containerNode:e,current:t,pendingChildren:null,callbackScheduled:!1,pendingLanes:0,suspendedLanes:0,expirationTimes:Array(8).fill(-1),callbackNode:null,callbackPriority:0,isHydrating:!1}}function nn(e,t){e.pendingChildren=t,e.callbackScheduled||(e.callbackScheduled=!0,oe(()=>{e.callbackScheduled=!1,sn(e)}))}function rn(e,t){e.pendingChildren=t,sn(e)}function an(e){let t=e;for(;t?.return;)t=t.return;if(t&&t.stateNode){let e=ln(t.stateNode);e&&nn(e,e.pendingChildren)}}function on(){wt(e=>{an(e)})}function sn(e){pt(),St(),on(),e.isHydrating&&(J=e,Y.clear());let t=e.current,n=Ve(t,{children:e.pendingChildren});cn(n),Nn(e,n),e.current=n,e.isHydrating&&(e.isHydrating=!1,J=null,Y.clear())}function cn(e){let t=e;for(;t!==null;)t=un(t)}function ln(e){for(let[,t]of en)if(t.containerNode===e)return t;return null}function un(e){if(dn(e),e.child!==null)return e.child;let t=e;for(;t!==null;){if(jn(t),t.sibling!==null)return t.sibling;t=t.return}return null}function dn(e){if(J!==null&&e.stateNode===null&&(e.tag===3||e.tag===4)){let t=An(e,J.containerNode);if(t)if(e.tag===3){let n=On(e,t);n&&(e.stateNode=n)}else{let n=kn(t);n&&(e.stateNode=n)}}switch(e.tag){case 2:fn(e);break;case 3:pn(e);break;case 4:break;case 0:mn(e);break;case 1:hn(e);break;case 5:e.type===h&&En++,gn(e);break;case 6:_n(e);break;case 8:vn(e);break;case 9:yn(e);break;default:gn(e);break}}function fn(e){let t=e.pendingProps.children;e.child=k(e,e.alternate?.child??null,t,0)}function pn(e){let t=e.pendingProps.children;e.child=k(e,e.alternate?.child??null,t,0)}function mn(e){Wt(),B(e);let t=e.type;mt(e,p(t)||t.name||`Anonymous`),En>0&&e.alternate===null&&(t(e.pendingProps),B(e));let n=t(e.pendingProps);e.dependencies=H(),B(null),Gt(),e.child=k(e,e.alternate?.child??null,n,0)}function hn(e){let t=e.type,n;if(e.stateNode===null?(n=new t(e.pendingProps),e.stateNode=n,n._fiber=e,n._enqueueUpdate=function(t){an(e),typeof t==`function`&&oe(t)}):(n=e.stateNode,n.props=e.pendingProps,n._fiber=e),n._pendingState.length>0||n._forceUpdate){let t=n.state;for(let e of n._pendingState)if(typeof e==`function`){let r=e(t,n.props);r!==null&&(t={...t,...r})}else t={...t,...e};n._pendingState=[];let r=n._forceUpdate||!n.shouldComponentUpdate||n.shouldComponentUpdate(n.props,t);if(n.state=t,n._forceUpdate=!1,!r){e.child=xn(e.alternate?.child??null,e);return}}let r=n.render();e.child=k(e,e.alternate?.child??null,r,0)}function gn(e){let t=e.pendingProps.children;e.child=k(e,e.alternate?.child??null,t,0)}function _n(e){let t=e.type,n=e.pendingProps.value;t._context&&(t._context._currentValue=n);let r=e.pendingProps.children;e.child=k(e,e.alternate?.child??null,r,0)}function vn(e){let{render:t}=e.type;Wt(),B(e);let n=t(e.pendingProps,e.ref);e.dependencies=H(),B(null),Gt(),e.child=k(e,e.alternate?.child??null,n,0)}function yn(e){let{type:t,compare:n}=e.type;if(e.alternate!==null){let t=e.alternate.memoizedProps,r=e.pendingProps;if(t!==null&&(n?n(t,r):bn(t,r))){e.child=xn(e.alternate.child,e);return}}Wt(),B(e);let r=t(e.pendingProps);e.dependencies=H(),B(null),Gt(),e.child=k(e,e.alternate?.child??null,r,0)}function bn(e,t){let n=Object.keys(e),r=Object.keys(t);if(n.length!==r.length)return!1;for(let r of n)if(r!==`children`&&!Object.is(e[r],t[r]))return!1;return!0}function xn(e,t){if(e===null)return null;let n=Sn(e,t),r=[],i=e.child;for(i&&r.push({source:i,cloneParent:n});r.length>0;){let{source:e,cloneParent:t}=r.pop(),n=e,i=null;for(;n!==null;){let e=Sn(n,t);i?i.sibling=e:t.child=e,n.child&&r.push({source:n.child,cloneParent:e}),i=e,n=n.sibling}}return n}function Sn(e,t){let n={...e,return:t,alternate:e,child:null,sibling:null,effectTag:0,pendingProps:e.memoizedProps??e.pendingProps};return e.alternate=n,n}var Cn=new Set(`svg.circle.clipPath.defs.ellipse.g.line.linearGradient.mask.path.pattern.polygon.polyline.radialGradient.rect.stop.text.tspan.use.image.symbol.foreignObject.desc.title.metadata.marker.filter.feBlend.feColorMatrix.feComponentTransfer.feComposite.feConvolveMatrix.feDiffuseLighting.feDisplacementMap.feFlood.feGaussianBlur.feImage.feMerge.feMergeNode.feMorphology.feOffset.feSpecularLighting.feTile.feTurbulence.animate.animateMotion.animateTransform.set`.split(`.`));function wn(e){return Cn.has(e)}function Tn(e){return e.tagName.includes(`-`)}var En=0,J=null,Y=new Map;function Dn(e){let t=e.firstChild;for(;t!==null;){if(t.nodeType===1||t.nodeType===3)return t;t=t.nextSibling}return null}function X(e){let t=e.nextSibling;for(;t!==null;){if(t.nodeType===1||t.nodeType===3)return t;t=t.nextSibling}return null}function On(e,t){let n=Y.get(t)??Dn(t);if(n===null)return null;for(;n!==null&&n.nodeType!==1;)n=X(n);if(n===null)return null;let r=n;return r.tagName.toLowerCase()===e.type.toLowerCase()?(Y.set(t,X(r)),r):null}function kn(e){let t=Y.get(e)??Dn(e);if(t===null)return null;for(;t!==null&&t.nodeType!==3;)t=X(t);if(t===null)return null;let n=t;return Y.set(e,X(n)),n}function An(e,t){let n=e.return;for(;n!==null;){if(n.tag===3&&n.stateNode)return n.stateNode;if(n.tag===2)return t;n=n.return}return t}function jn(e){switch(e.tag){case 3:if(e.stateNode===null){let t=e.type,n=wn(t)?document.createElementNS(`http://www.w3.org/2000/svg`,t):document.createElement(t);Z(n,{},e.pendingProps),$t(t,e.pendingProps);{let t=e.return;for(;t!==null;){if(t.tag===0||t.tag===1||t.tag===8||t.tag===9){if(t.child===e){let e=ie(p(t.type));e&&n.setAttribute(`data-cid`,e)}break}if(t.tag===3||t.tag===2)break;t=t.return}}e.stateNode=n,Mn(n,e)}else if(J!==null&&e.alternate===null)Z(e.stateNode,{},e.pendingProps),e.effectTag=0;else{let t=e.stateNode;e.alternate&&Z(t,e.alternate.memoizedProps||{},e.pendingProps)}e.memoizedProps=e.pendingProps;break;case 4:{let t=String(e.pendingProps.text);if(e.stateNode===null)e.stateNode=document.createTextNode(t);else if(J!==null&&e.alternate===null){let n=e.stateNode;n.nodeValue!==t&&(n.nodeValue=t),e.effectTag=0}else e.stateNode.nodeValue=t;e.memoizedProps=e.pendingProps;break}case 2:case 0:case 1:case 5:case 6:case 8:case 9:e.type===h&&En--,e.memoizedProps=e.pendingProps;break}}function Mn(e,t){let n=t.child;for(;n!==null;){if(n.tag===3||n.tag===4)n.stateNode&&e.appendChild(n.stateNode);else if(n.child!==null){n.child.return=n,n=n.child;continue}if(n===t)return;for(;n.sibling===null;){if(n.return===null||n.return===t)return;n=n.return}n.sibling.return=n.return,n=n.sibling}}function Nn(e,t){_t(),Pn(t),Ln(t,e.containerNode),Vn(t),qt(e)}function Pn(e){let t=e;for(;t!==null;){if(t.updateQueue&&Array.isArray(t.updateQueue))for(let e of t.updateQueue){let t=e;t.effectTag===4&&Fn(t)}if(t.child!==null){t=t.child;continue}for(;t!==null;){if(t===e){t=null;break}if(t.sibling!==null){t=t.sibling;break}t=t.return}}}function Fn(e){let t=e.return;for(;t!==null&&!(t.tag===3||t.tag===2);)t=t.return;let n=t?.stateNode;n&&(In(e,n),Un(e))}function In(e,t){let n=e;for(;n!==null;){if(n.tag===3||n.tag===4)n.stateNode&&t.contains(n.stateNode)&&t.removeChild(n.stateNode);else if(n.child!==null){n=n.child;continue}for(;n!==null;){if(n===e){n=null;break}if(n.sibling!==null){n=n.sibling;break}n=n.return}}}function Ln(e,t){let n=e;for(;n!==null;){if(n.tag!==2&&(n.tag===3||n.tag===4)&&n.effectTag&1){let e=Rn(n,t);if(e&&n.stateNode){let t=zn(n);t?e.insertBefore(n.stateNode,t):e.appendChild(n.stateNode)}}if(n.child!==null){n=n.child;continue}for(;n!==null;){if(n===e){n=null;break}if(n.sibling!==null){n=n.sibling;break}n=n.return}}}function Rn(e,t){let n=e.return;for(;n!==null;){if(n.tag===3)return n.stateNode;if(n.tag===2)return t;n=n.return}return null}function zn(e){let t=e;outer:for(;;){for(;t.sibling===null;){if(t.return===null||Bn(t.return))return null;t=t.return}for(t=t.sibling;t.tag!==3&&t.tag!==4;){if(t.effectTag&1||t.child===null)continue outer;t=t.child}if(!(t.effectTag&1))return t.stateNode}}function Bn(e){return e.tag===3||e.tag===2}function Vn(e){let t=e;for(;t!==null;){if(t.ref&&t.stateNode&&(typeof t.ref==`function`?t.ref(t.stateNode):typeof t.ref==`object`&&(t.ref.current=t.stateNode)),t.tag===0||t.tag===8||t.tag===9){let e=t.dependencies;if(e){let n=(typeof t.type==`function`?p(t.type)||t.type.name:``)||`Anonymous`;vt(t,n),Hn(e)}}if(t.tag===1&&t.stateNode){let e=t.stateNode;if(t.alternate===null)e.componentDidMount?.();else{let n=t.alternate.memoizedProps||{},r=t.alternate.memoizedState;e.componentDidUpdate?.(n,r)}}if(t.child!==null){t=t.child;continue}for(;t!==null;){if(t===e){t=null;break}if(t.sibling!==null){t=t.sibling;break}t=t.return}}}function Hn(e){let t=e;for(;t!==null;){if(t.tag&1){t.destroy&&t.destroy();let e=t.create();t.destroy=typeof e==`function`?e:null}t=t.next}}function Un(e){let t=e;for(;t!==null;){if(t.tag===0||t.tag===8){let e=t.dependencies;if(e){let t=e;for(;t!==null;)t.destroy&&t.destroy(),t=t.next}}if(t.tag===1&&t.stateNode&&t.stateNode.componentWillUnmount?.(),t.child!==null){t=t.child;continue}for(;t!==null;){if(t===e){t=null;break}if(t.sibling!==null){t=t.sibling;break}t=t.return}}}var Wn=/^on[A-Z]/;function Z(e,t,n){for(let r in t)if(!(r===`children`||r===`key`||r===`ref`)&&!(r in n))if(Wn.test(r)){let n=r.slice(2).toLowerCase();e.removeEventListener(n,t[r])}else r===`style`?e.removeAttribute(`style`):r===`className`?e.removeAttribute(`class`):r===`htmlFor`?e.removeAttribute(`for`):r===`dangerouslySetInnerHTML`?e.innerHTML=``:e.removeAttribute(r);for(let r in n){if(r===`children`||r===`key`||r===`ref`)continue;let i=n[r];if(Wn.test(r)){let n=r.slice(2).toLowerCase();t[r]!==i&&(t[r]&&e.removeEventListener(n,t[r]),i&&e.addEventListener(n,i))}else if(r===`style`){if(typeof i==`object`&&i){let t=i;for(let n in t)e.style[n]=t[n]??``}}else r===`className`?e.setAttribute(`class`,String(i)):r===`htmlFor`?e.setAttribute(`for`,String(i)):r===`dangerouslySetInnerHTML`?e.innerHTML=i.__html:r===`value`&&(e.tagName===`INPUT`||e.tagName===`TEXTAREA`||e.tagName===`SELECT`)?e.value=String(i??``):r===`checked`&&e.tagName===`INPUT`?e.checked=!!i:typeof i==`boolean`?i?e.setAttribute(r,``):e.removeAttribute(r):i==null?e.removeAttribute(r):Tn(e)&&typeof i!=`string`&&r in e?e[r]=i:e.setAttribute(r,String(i))}}function Gn(e,t){if(!e||!(e instanceof Element)&&!(e instanceof DocumentFragment))throw Error(`createRoot: container must be a DOM element or DocumentFragment`);let n=tn(e);en.set(e,n);let r=!1;return{render(e){r=!0,rn(n,e)},unmount(){r&&(r=!1,rn(n,null),en.delete(e))}}}function Q(e,t,n){let r=n===void 0?t.key===void 0?null:t.key:n,{key:i,...a}=t;return _(e,r===null?a:{...a,key:r})}function Kn(){let{pathname:e}=Re();return Q(`nav`,{class:`nav`,children:Q(`div`,{class:`nav-inner`,children:[Q(T,{to:`/`,class:`nav-brand`,children:[Q(`img`,{src:`/docs/img/logo.png`,alt:`datascience logo`,width:`28`,height:`28`}),`datascience`]}),Q(`div`,{class:`nav-links`,children:[Q(T,{to:`/`,class:e===`/`?`active`:``,children:`Home`}),Q(T,{to:`/docs`,class:e===`/docs`?`active`:``,children:`Docs`}),Q(T,{to:`/tutorials`,class:e===`/tutorials`?`active`:``,children:`Tutorials`}),Q(T,{to:`/api`,class:e===`/api`?`active`:``,children:`API`}),Q(T,{to:`/libraries`,class:e===`/libraries`?`active`:``,children:`Libraries`})]})]})})}function qn(){return Q(`footer`,{class:`footer`,children:Q(`div`,{class:`footer-inner`,children:[Q(`span`,{children:`v0.1.3 | MIT License © 2026 Asymmetric Effort, LLC`}),Q(`div`,{class:`footer-links`,children:[Q(`a`,{href:`https://github.com/asymmetric-effort/datascience`,target:`_blank`,rel:`noopener noreferrer`,children:`GitHub`}),Q(`a`,{href:`https://github.com/asymmetric-effort/datascience/blob/main/SECURITY.md`,target:`_blank`,rel:`noopener noreferrer`,children:`Security`}),Q(`a`,{href:`https://github.com/asymmetric-effort/datascience/blob/main/CONTRIBUTING.md`,target:`_blank`,rel:`noopener noreferrer`,children:`Contributing`})]})]})})}function Jn(){return y({title:`datascience — Data Science and Machine Learning in Go`,description:`A comprehensive pure Go library for data science and machine learning. PGMs, TensorFlow/deep learning, BLAS, financial modeling, and more.`,canonical:`https://datascience.asymmetric-effort.com/`,og:{title:`datascience — Data Science and Machine Learning in Go`,description:`A comprehensive pure Go library for data science and machine learning.`,url:`https://datascience.asymmetric-effort.com/`}}),Q(`div`,{class:`page`,children:[Q(`section`,{class:`hero`,children:[Q(`img`,{src:`/docs/img/logo.png`,alt:`datascience logo`,class:`hero-logo`}),Q(`h1`,{children:`datascience`}),Q(`p`,{class:`hero-subtitle`,children:`Data Science and Machine Learning in Go`}),Q(`p`,{class:`hero-description`,children:`A comprehensive pure Go library for probabilistic graphical models, TensorFlow-compatible deep learning, BLAS linear algebra, financial modeling, and more. Built entirely in Go with near-zero dependencies.`}),Q(`div`,{class:`badges`,children:[Q(`span`,{class:`badge`,children:`Zero Dependencies`}),Q(`span`,{class:`badge`,children:`Go 1.26+`}),Q(`span`,{class:`badge`,children:`MIT License`}),Q(`span`,{class:`badge`,children:`~5,000 Tests`}),Q(`span`,{class:`badge`,children:`392 Cross-Validation Fixtures`}),Q(`span`,{class:`badge`,children:`PGMs`}),Q(`span`,{class:`badge`,children:`TensorFlow / Deep Learning`}),Q(`span`,{class:`badge`,children:`BLAS`}),Q(`span`,{class:`badge`,children:`Financial Modeling`}),Q(`span`,{class:`badge`,children:`40 Built-in Datasets`})]}),Q(`div`,{class:`hero-actions`,children:[Q(T,{to:`/docs`,class:`btn btn-primary`,children:`Get Started`}),Q(`a`,{href:`https://github.com/asymmetric-effort/datascience`,target:`_blank`,rel:`noopener noreferrer`,class:`btn btn-secondary`,children:`GitHub`})]})]}),Q(`section`,{class:`section`,children:[Q(`h2`,{children:`Installation`}),Q(`pre`,{children:Q(`code`,{children:`go get github.com/asymmetric-effort/datascience`})}),Q(`p`,{children:`Requires Go 1.21 or later. No C dependencies, no cgo, no third-party modules.`})]}),Q(`section`,{class:`section`,children:[Q(`h2`,{children:`Quick Start`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
)

func main() {
    // Option 1: Load a built-in example model
    asia, _ := example_models.Get("asia")
    fmt.Println("Asia model:", len(asia.Nodes()), "nodes")

    // Option 2: Build a Bayesian Network from scratch
    bn := models.NewBayesianNetwork()
    bn.AddNode("Rain")
    bn.AddNode("Sprinkler")
    bn.AddNode("WetGrass")
    bn.AddEdge("Rain", "WetGrass")
    bn.AddEdge("Sprinkler", "WetGrass")

    // Add CPDs
    bn.SetStates("Rain", []string{"no", "yes"})
    bn.SetStates("Sprinkler", []string{"off", "on"})
    bn.SetStates("WetGrass", []string{"dry", "wet"})

    bn.SetCPD("Rain", factors.NewTabularCPD(
        "Rain", 2, []float64{0.8, 0.2}, nil, nil,
    ))

    // Run inference with Variable Elimination
    facs, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    result, _ := ve.Query(
        []string{"WetGrass"},
        map[string]int{"Rain": 1},
    )
    fmt.Println(result)
}`})}),Q(`p`,{children:[`See the `,Q(T,{to:`/docs`,children:`documentation`}),` for complete getting-started guides, or jump to `,Q(T,{to:`/tutorials`,children:`tutorials`}),` for step-by-step walkthroughs.`]})]}),Q(`section`,{class:`section`,children:[Q(`h2`,{children:`Features`}),Q(`div`,{class:`features-grid`,children:[Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Probabilistic Graphical Models`}),Q(`p`,{children:`13 model types including BayesianNetwork, MarkovNetwork, DynamicBN, NaiveBayes, SEM, FactorGraph, JunctionTree, and more. 7 inference algorithms, 15+ learning algorithms.`})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`TensorFlow / Deep Learning`}),Q(`p`,{children:[`TensorFlow-compatible deep learning via the `,Q(`code`,{children:`lib/tensorflow`}),` package. Neural network construction, training, and inference in pure Go.`]})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`BLAS Linear Algebra`}),Q(`p`,{children:[`High-performance BLAS routines for matrix multiplication, decomposition, eigenvalue computation, and other linear algebra operations via `,Q(`code`,{children:`lib/numgo`}),`.`]})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Financial Modeling`}),Q(`p`,{children:`Tools for quantitative finance including time series analysis, risk modeling, portfolio optimization, and statistical forecasting.`})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Statistical Computing`}),Q(`p`,{children:[`Comprehensive statistical distributions, hypothesis tests, optimization routines, and special functions via `,Q(`code`,{children:`lib/scigo`}),` (scipy equivalent).`]})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Graph Algorithms`}),Q(`p`,{children:[`Directed and undirected graphs with topological sort, d-separation, connected components, shortest paths, and more via `,Q(`code`,{children:`lib/graphgo`}),` (networkx equivalent).`]})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Tabular Data`}),Q(`p`,{children:[`DataFrames, Series, CSV/Parquet/Excel I/O, filtering, groupby, and merge operations via `,Q(`code`,{children:`lib/tabgo`}),` (pandas equivalent).`]})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`GPU Compute Backend`}),Q(`p`,{children:[`Optional GPU acceleration via the `,Q(`code`,{children:`lib/gpu`}),` package for compute-intensive operations on large networks and deep learning models.`]})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Causal Inference`}),Q(`p`,{children:`Do-calculus interventional queries, back-door and front-door identification, ATE estimation with DoubleML, naive adjustment, and IV regression.`})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`LLM Integration`}),Q(`p`,{children:`Expert-in-the-loop structure learning with LLM client support for AI-assisted model construction and knowledge elicitation.`})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`40 Built-in Datasets`}),Q(`p`,{children:`Ready-to-use CSV datasets for structure learning, parameter estimation, benchmarking, and general machine learning tasks.`})]}),Q(`div`,{class:`feature-card`,children:[Q(`h3`,{children:`Near-Zero Dependencies`}),Q(`p`,{children:`Built entirely in Go with custom implementations of numpy (numgo), scipy (scigo), networkx (graphgo), and pandas (tabgo). No cgo, no external libraries.`})]})]})]}),Q(`section`,{class:`section`,children:[Q(`h2`,{children:`Why datascience?`}),Q(`p`,{children:`datascience brings comprehensive data science and machine learning capabilities to the Go ecosystem. Whether you need probabilistic graphical models (inspired by pgmpy), deep learning (TensorFlow-compatible), BLAS linear algebra, or financial modeling -- this library provides a unified, pure Go solution.`}),Q(`p`,{children:[`Every numerical primitive -- linear algebra, statistical distributions, graph algorithms, tabular data processing -- is implemented from scratch in pure Go. This means `,Q(`code`,{children:`go build`}),` just works, cross-compilation just works, and static binaries just work.`]}),Q(`p`,{children:[`The layered architecture means you can use individual packages independently. Need just matrix math? Import `,Q(`code`,{children:`lib/numgo`}),`. Just graphs? Import `,Q(`code`,{children:`lib/graphgo`}),`. Just PGMs? Import `,Q(`code`,{children:`lib/pgm`}),`. Just deep learning? Import `,Q(`code`,{children:`lib/tensorflow`}),`.`]})]})]})}function $({to:e,children:t,class:n}){return Q(`a`,{href:`javascript:void(0)`,onClick:t=>{t.preventDefault();let n=document.getElementById(e);n&&n.scrollIntoView({behavior:`smooth`,block:`start`})},class:n||``,children:t})}function Yn(){return y({title:`Documentation — datascience`,description:`Comprehensive documentation for datascience, a pure Go library for data science and machine learning. PGMs, TensorFlow, BLAS, financial modeling, and more.`,canonical:`https://datascience.asymmetric-effort.com/#/docs`}),Q(`div`,{class:`page`,children:[Q(`h1`,{children:`Documentation`}),Q(`nav`,{class:`page-toc`,children:[Q(`strong`,{children:`On this page:`}),` `,Q($,{to:`overview`,children:`Overview`}),` | `,Q($,{to:`installation`,children:`Installation`}),` | `,Q($,{to:`getting-started`,children:`Getting Started`}),` | `,Q($,{to:`architecture`,children:`Architecture`}),` | `,Q($,{to:`library-packages`,children:`Library Packages`}),` | `,Q($,{to:`pgm-packages`,children:`PGM Packages`}),` | `,Q($,{to:`file-formats`,children:`File Formats`}),` | `,Q($,{to:`example-models`,children:`Example Models`}),` | `,Q($,{to:`datasets`,children:`Datasets`}),` | `,Q($,{to:`testing`,children:`Testing`}),` | `,Q($,{to:`configuration`,children:`Configuration`}),` | `,Q($,{to:`contributing`,children:`Contributing`})]}),Q(`section`,{class:`section`,id:`overview`,children:[Q(`h2`,{children:`Overview`}),Q(`p`,{children:[Q(`strong`,{children:`datascience`}),` is a comprehensive pure Go library for data science and machine learning. It provides tools for probabilistic graphical models (PGMs), TensorFlow-compatible deep learning, BLAS linear algebra, financial modeling, statistical computing, and more.`]}),Q(`p`,{children:[`All numerical primitives -- linear algebra, statistical distributions, graph algorithms, and tabular data processing -- are implemented from scratch in pure Go. There are no C bindings, no cgo, no third-party module dependencies. The result is a library that compiles to a single static binary with `,Q(`code`,{children:`go build`}),`, cross-compiles trivially, and deploys without runtime dependencies.`]}),Q(`p`,{children:`The library includes approximately 5,000 tests and 392 cross-validation fixtures covering inference, learning, sampling, serialization, and cross-validation across all packages.`}),Q(`h3`,{children:`Key Capabilities`}),Q(`ul`,{children:[Q(`li`,{children:[Q(`strong`,{children:`Probabilistic Graphical Models`}),`: 13 model types, 7 inference algorithms, 15+ learning algorithms, 16 CI tests, 13 scoring functions (via `,Q(`code`,{children:`lib/pgm`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`TensorFlow / Deep Learning`}),`: Neural network construction, training, and inference (via `,Q(`code`,{children:`lib/tensorflow`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`BLAS Linear Algebra`}),`: N-dimensional arrays, matrices, vectors, decompositions (via `,Q(`code`,{children:`lib/numgo`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`Statistical Computing`}),`: Distributions, optimization, hypothesis tests, special functions (via `,Q(`code`,{children:`lib/scigo`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`Graph Algorithms`}),`: Directed/undirected graphs, topological sort, d-separation, connected components (via `,Q(`code`,{children:`lib/graphgo`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`Tabular Data`}),`: DataFrames, Series, CSV/Parquet/Excel I/O (via `,Q(`code`,{children:`lib/tabgo`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`GPU Acceleration`}),`: Optional GPU compute backend for large-scale operations (via `,Q(`code`,{children:`lib/gpu`}),`)`]}),Q(`li`,{children:[Q(`strong`,{children:`10 file formats`}),`: BIF, XMLBIF, NET, UAI, XDSL, PomdpX, XBN, CSV, JSON, XML`]}),Q(`li`,{children:[Q(`strong`,{children:`25 built-in example models`}),` and `,Q(`strong`,{children:`40 built-in datasets`})]}),Q(`li`,{children:[Q(`strong`,{children:`LLM integration`}),` for expert-in-the-loop structure learning`]})]}),Q(`h3`,{children:`Relationship to pgmpy`}),Q(`p`,{children:[`The PGM module (`,Q(`code`,{children:`lib/pgm`}),`) is inspired by `,Q(`a`,{href:`https://pgmpy.org`,target:`_blank`,rel:`noopener noreferrer`,children:`pgmpy`}),` and follows similar API patterns where possible. If you have used pgmpy in Python, the concepts and workflows will feel familiar. However, datascience is not a direct port -- it is a ground-up reimplementation in Go with its own design decisions, particularly around the near-zero-dependency philosophy.`]}),Q(`p`,{children:[`Where pgmpy relies on numpy, scipy, networkx, and pandas, datascience provides its own equivalents:`,Q(`code`,{children:`numgo`}),`, `,Q(`code`,{children:`scigo`}),`, `,Q(`code`,{children:`graphgo`}),`, and `,Q(`code`,{children:`tabgo`}),`. These are general-purpose libraries that can be used independently of the PGM layer.`]})]}),Q(`section`,{class:`section`,id:`installation`,children:[Q(`h2`,{children:`Installation`}),Q(`h3`,{children:`Go Module`}),Q(`pre`,{children:Q(`code`,{children:`go get github.com/asymmetric-effort/datascience`})}),Q(`p`,{children:[`Requires `,Q(`strong`,{children:`Go 1.21`}),` or later. No C compiler needed, no system libraries needed. Works on Linux, macOS, and Windows. Cross-compilation works out of the box.`]}),Q(`h3`,{children:`From Source`}),Q(`pre`,{children:Q(`code`,{children:`git clone https://github.com/asymmetric-effort/datascience.git
cd datascience
go test ./...`})}),Q(`h3`,{children:`Verify`}),Q(`pre`,{children:Q(`code`,{children:`# Run all tests to verify the installation
go test ./...

# Run a quick sanity check
go test ./lib/pgm/models/... -run TestBayesianNetwork -v`})})]}),Q(`section`,{class:`section`,id:`getting-started`,children:[Q(`h2`,{children:`Getting Started`}),Q(`h3`,{children:`Your First Bayesian Network`}),Q(`p`,{children:`A Bayesian network is a directed acyclic graph (DAG) where nodes represent random variables and edges represent conditional dependencies. Here is a complete, runnable program that creates the classic "wet grass" network:`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
)

func main() {
    // Create an empty Bayesian network
    bn := models.NewBayesianNetwork()

    // Add nodes (random variables)
    bn.AddNode("Cloudy")
    bn.AddNode("Rain")
    bn.AddNode("Sprinkler")
    bn.AddNode("WetGrass")

    // Add directed edges (causal relationships)
    bn.AddEdge("Cloudy", "Rain")
    bn.AddEdge("Cloudy", "Sprinkler")
    bn.AddEdge("Rain", "WetGrass")
    bn.AddEdge("Sprinkler", "WetGrass")

    // Define states for each variable
    bn.SetStates("Cloudy", []string{"clear", "cloudy"})
    bn.SetStates("Rain", []string{"no", "yes"})
    bn.SetStates("Sprinkler", []string{"off", "on"})
    bn.SetStates("WetGrass", []string{"dry", "wet"})

    // P(Cloudy) -- root node, marginal distribution
    bn.SetCPD("Cloudy", factors.NewTabularCPD(
        "Cloudy", 2,
        []float64{0.5, 0.5},
        nil, nil,
    ))

    // P(Rain | Cloudy) -- conditional distribution
    bn.SetCPD("Rain", factors.NewTabularCPD(
        "Rain", 2,
        []float64{0.8, 0.2, 0.2, 0.8},
        []string{"Cloudy"}, []int{2},
    ))

    // P(Sprinkler | Cloudy)
    bn.SetCPD("Sprinkler", factors.NewTabularCPD(
        "Sprinkler", 2,
        []float64{0.5, 0.9, 0.5, 0.1},
        []string{"Cloudy"}, []int{2},
    ))

    // P(WetGrass | Sprinkler, Rain)
    bn.SetCPD("WetGrass", factors.NewTabularCPD(
        "WetGrass", 2,
        []float64{
            1.0, 0.1, 0.1, 0.01,
            0.0, 0.9, 0.9, 0.99,
        },
        []string{"Sprinkler", "Rain"}, []int{2, 2},
    ))

    // Validate: checks DAG, CPD dimensions, probability sums
    if err := bn.CheckModel(); err != nil {
        log.Fatal("Model error:", err)
    }
    fmt.Println("Model is valid!")
    fmt.Printf("Nodes: %d, Edges: %d\\n", len(bn.Nodes()), len(bn.Edges()))
}`})}),Q(`h3`,{children:`Your First Query`}),Q(`p`,{children:`Once you have a valid model, convert its CPDs to Markov factors and use Variable Elimination to compute posterior probabilities:`}),Q(`pre`,{children:Q(`code`,{children:`import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

func main() {
    // Load a built-in model (13 models have full CPDs)
    bn, _ := example_models.Get("asia")

    // Convert CPDs to Markov factors
    facs, err := bn.ToMarkovFactors()
    if err != nil {
        log.Fatal(err)
    }

    // Create a Variable Elimination engine
    ve := inference.NewVariableElimination(facs)

    // Posterior query: P(Dyspnea | Smoker=yes)
    result, err := ve.Query(
        []string{"Dyspnea"},           // query variables
        map[string]int{"Smoker": 1},   // evidence (1 = "yes")
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Dyspnea | Smoker=yes):", result.Values().Data())

    // MAP query: most likely assignment
    assignment, err := ve.MAP(
        []string{"Lung", "Bronc"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("MAP(Lung, Bronc | Smoker=yes):", assignment)
}`})}),Q(`h3`,{children:`Your First Structure Learning`}),Q(`p`,{children:`When you have observational data but no known structure, use structure learning to discover the DAG:`}),Q(`pre`,{children:Q(`code`,{children:`import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/learning"
    "github.com/asymmetric-effort/datascience/lib/pgm/sampling"
    "github.com/asymmetric-effort/datascience/lib/pgm/structure_score"
)

func main() {
    // Generate training data from a known model
    bn, _ := example_models.Get("asia")
    bms, _ := sampling.NewBayesianModelSampling(bn, 42)
    data, _ := bms.ForwardSample(5000)

    // Learn structure using hill-climbing with BIC scoring
    scorer := structure_score.NewBIC()
    hc := learning.NewHillClimbSearch(data, scorer.LocalScore)
    learnedBN, err := hc.Estimate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Learned: %d nodes, %d edges\\n",
        len(learnedBN.Nodes()), len(learnedBN.Edges()))

    // Fit parameters to the learned structure
    mle := learning.NewMLE(learnedBN, data)
    if err := mle.Estimate(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Parameters fitted successfully")
}`})})]}),Q(`section`,{class:`section`,id:`architecture`,children:[Q(`h2`,{children:`Architecture`}),Q(`h3`,{children:`Project Layout`}),Q(`pre`,{children:Q(`code`,{children:`datascience/
  lib/                     # All library modules
    numgo/                 # numpy equivalent: NDArray, Matrix, Vector, BLAS
    scigo/                 # scipy equivalent: distributions, optimization, special functions
    graphgo/               # networkx equivalent: DiGraph, PDAG, graph algorithms
    tabgo/                 # pandas equivalent: DataFrame, Series, CSV/Parquet I/O
    gpu/                   # GPU compute backend for large-scale operations
    pgm/                   # Probabilistic Graphical Models (pgmpy-inspired)
      base/                # DAG, PDAG, UndirectedGraph, ADMG, MAG, SimpleCausalModel
      models/              # 13 probabilistic model types
      factors/             # Factor representations: TabularCPD, DiscreteFactor, etc.
      inference/           # 7 inference algorithms (VE, BP, MPLP, Causal, etc.)
      sampling/            # Forward, rejection, likelihood-weighted, Gibbs sampling
      learning/            # 15+ learning algorithms (parameter + structure)
      ci_tests/            # 16 conditional independence tests
      structure_score/     # 13 scoring functions for structure learning
      identification/      # Causal effect identification (back-door, front-door)
      prediction/          # DoubleML, naive adjustment, IV regression
      metrics/             # SHD, confusion matrices, correlation, Fisher's C
      independencies/      # Independence assertion representations
      readwrite/           # 10 file format readers and writers
      config/              # Global configuration
      utils/               # Shared parsing, optimization, compatibility utilities
    tensorflow/            # TensorFlow-compatible deep learning
  example_models/          # 25 built-in Bayesian networks
  examples/
    datasets/              # 40 built-in CSV datasets
  website/                 # Project website (this site)
  docs/                    # Additional documentation`})}),Q(`h3`,{children:`Dependency Flow`}),Q(`p`,{children:`The dependency flow is strictly layered. Each layer depends only on layers below it:`}),Q(`pre`,{children:Q(`code`,{children:`Layer 3: lib/pgm, lib/tensorflow  (domain: PGMs, deep learning, etc.)
         |
Layer 2: lib/numgo, lib/scigo,    (primitives: arrays, stats, graphs, tables, GPU)
         lib/graphgo, lib/tabgo,
         lib/gpu
         |
Layer 1: Go standard library       (only dependency)`})}),Q(`p`,{children:[`The `,Q(`code`,{children:`lib/`}),` packages are general-purpose and can be used independently. For example, you could use `,Q(`code`,{children:`numgo`}),` for matrix math or `,Q(`code`,{children:`graphgo`}),` for graph algorithms without importing any PGM-specific code.`]}),Q(`p`,{children:[`The `,Q(`code`,{children:`lib/pgm`}),` packages build on the primitive libraries to implement PGM-specific functionality. They also depend on each other -- for example, `,Q(`code`,{children:`inference`}),` depends on `,Q(`code`,{children:`factors`}),`, and `,Q(`code`,{children:`learning`}),` depends on `,Q(`code`,{children:`structure_score`}),` and `,Q(`code`,{children:`ci_tests`}),`.`]}),Q(`p`,{children:[`The `,Q(`code`,{children:`lib/tensorflow`}),` package provides TensorFlow-compatible deep learning functionality, building on `,Q(`code`,{children:`numgo`}),` for tensor operations and `,Q(`code`,{children:`gpu`}),` for acceleration.`]}),Q(`h3`,{children:`Design Principles`}),Q(`ul`,{children:[Q(`li`,{children:[Q(`strong`,{children:`Near-zero dependencies:`}),` The entire library compiles with only the Go standard library. No cgo, no system libraries, no third-party modules.`]}),Q(`li`,{children:[Q(`strong`,{children:`pgmpy compatibility:`}),` PGM API patterns follow pgmpy where practical, making it easier for Python practitioners to transition to Go.`]}),Q(`li`,{children:[Q(`strong`,{children:`Cross-validation:`}),` 392 test fixtures validate results against known-correct outputs, ensuring numerical accuracy across inference, learning, and sampling.`]}),Q(`li`,{children:[Q(`strong`,{children:`Layered architecture:`}),` Primitive libraries (numgo, scigo, graphgo, tabgo) are reusable beyond PGMs. Domain libraries (pgm, tensorflow) build on these without polluting them.`]})]})]}),Q(`section`,{class:`section`,id:`library-packages`,children:[Q(`h2`,{children:`Library Packages (lib/)`}),Q(`p`,{children:`These packages replace common Python scientific computing libraries with pure Go implementations. They are general-purpose and can be imported independently of the PGM or deep learning layers.`}),Q(`h3`,{children:`numgo -- numpy equivalent`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/numgo`})]}),Q(`p`,{children:`N-dimensional arrays, linear algebra, BLAS routines, matrix operations, broadcasting, and element-wise arithmetic.`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NDArray`})}),Q(`td`,{children:`N-dimensional array with shape, stride, and element-wise operations (add, multiply, etc.)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Matrix`})}),Q(`td`,{children:`2D matrix with multiply, transpose, inverse, determinant, eigenvalues`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Vector`})}),Q(`td`,{children:`1D vector with dot product, norm, element-wise operations`})]})]})]}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/numgo"

// Create a matrix
m := numgo.NewMatrix(3, 3)
m.Set(0, 0, 1.0)
m.Set(1, 1, 2.0)
m.Set(2, 2, 3.0)

// Matrix operations
det := m.Det()
inv := m.Inverse()
transposed := m.Transpose()
product := m.Multiply(inv) // should be identity

// NDArray operations
arr := numgo.NewNDArray([]int{2, 3, 4}) // 2x3x4 array
arr.Fill(1.0)
sum := arr.Sum()

// Vector operations
v1 := numgo.NewVector([]float64{1, 2, 3})
v2 := numgo.NewVector([]float64{4, 5, 6})
dot := v1.Dot(v2)
norm := v1.Norm()`})}),Q(`h3`,{children:`scigo -- scipy equivalent`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/scigo`})]}),Q(`p`,{children:`Statistical distributions, optimization routines, special functions, and hypothesis tests.`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Category`}),Q(`th`,{children:`Key Types / Functions`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`Distributions`}),Q(`td`,{children:[Q(`code`,{children:`Normal`}),`, `,Q(`code`,{children:`ChiSquared`}),`, `,Q(`code`,{children:`Beta`}),`, `,Q(`code`,{children:`Gamma`}),`, `,Q(`code`,{children:`StudentT`}),`, `,Q(`code`,{children:`Uniform`}),`, `,Q(`code`,{children:`Exponential`})]})]}),Q(`tr`,{children:[Q(`td`,{children:`Optimization`}),Q(`td`,{children:[Q(`code`,{children:`Minimize`}),`, `,Q(`code`,{children:`GradientDescent`}),`, `,Q(`code`,{children:`NewtonMethod`})]})]}),Q(`tr`,{children:[Q(`td`,{children:`Statistics`}),Q(`td`,{children:`Hypothesis tests, p-value computation, quantile functions, CDF/PDF/PPF`})]}),Q(`tr`,{children:[Q(`td`,{children:`Special Functions`}),Q(`td`,{children:`Gamma function, beta function, incomplete gamma, digamma`})]})]})]}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/scigo"

// Normal distribution
n := scigo.NewNormal(0, 1) // mean=0, std=1
pdf := n.PDF(1.96)
cdf := n.CDF(1.96)       // ~0.975
ppf := n.PPF(0.975)      // ~1.96

// Chi-squared distribution (used in CI tests)
chi2 := scigo.NewChiSquared(5) // 5 degrees of freedom
pValue := 1.0 - chi2.CDF(11.07)

// Optimization
result := scigo.Minimize(func(x float64) float64 {
    return (x - 3) * (x - 3)
}, 0.0, 10.0)`})}),Q(`h3`,{children:`graphgo -- networkx equivalent`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/graphgo`})]}),Q(`p`,{children:`Directed and undirected graphs with a full suite of graph algorithms.`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiGraph`})}),Q(`td`,{children:`Directed graph with adjacency operations, successors, predecessors`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Graph`})}),Q(`td`,{children:`Undirected graph with neighbors, degree, connected components`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`PDAG`})}),Q(`td`,{children:`Partially directed acyclic graph (for equivalence classes)`})]})]})]}),Q(`p`,{children:`Algorithms: topological sort, d-separation, moral graph, triangulation, maximum cardinality search, clique finding, connected components, shortest paths, cycle detection.`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/graphgo"

// Create a directed graph
g := graphgo.NewDiGraph()
g.AddNode("A")
g.AddNode("B")
g.AddNode("C")
g.AddEdge("A", "B")
g.AddEdge("B", "C")

// Graph queries
parents := g.Predecessors("C")  // ["B"]
children := g.Successors("A")   // ["B"]
sorted := g.TopologicalSort()   // ["A", "B", "C"]

// Undirected graph
ug := graphgo.NewGraph()
ug.AddEdge("X", "Y")
ug.AddEdge("Y", "Z")
neighbors := ug.Neighbors("Y")  // ["X", "Z"]`})}),Q(`h3`,{children:`tabgo -- pandas equivalent`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/tabgo`})]}),Q(`p`,{children:`Tabular data with named columns, row filtering, groupby, and file I/O.`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DataFrame`})}),Q(`td`,{children:`Tabular data with named columns, row filtering, groupby, merge`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Series`})}),Q(`td`,{children:`Single column with value counts, unique values, statistical summaries`})]})]})]}),Q(`p`,{children:[`I/O: `,Q(`code`,{children:`ReadCSV`}),`, `,Q(`code`,{children:`WriteCSV`}),`, `,Q(`code`,{children:`ReadParquet`}),`, `,Q(`code`,{children:`WriteParquet`}),`, `,Q(`code`,{children:`ReadXLSX`}),`.`]}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/tabgo"

// Read CSV data
df, err := tabgo.ReadCSV("observations.csv")
fmt.Printf("Rows: %d, Columns: %d\\n", df.NRows(), len(df.Columns()))

// Access a column as a Series
col := df.Column("Temperature")
fmt.Println("Unique values:", col.Unique())
fmt.Println("Value counts:", col.ValueCounts())

// Filter rows
filtered := df.Filter(func(row map[string]interface{}) bool {
    return row["Temperature"].(int) > 70
})

// Write CSV
tabgo.WriteCSV(filtered, "warm_days.csv")`})}),Q(`h3`,{children:`gpu -- GPU Compute Backend`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/gpu`})]}),Q(`p`,{children:`Optional GPU acceleration for compute-intensive operations on large networks. Provides GPU-backed matrix operations and factor computations that can significantly speed up inference, learning, and deep learning training.`}),Q(`h3`,{children:`tensorflow -- Deep Learning`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/tensorflow`})]}),Q(`p`,{children:[`TensorFlow-compatible deep learning in pure Go. Provides neural network construction, training, and inference capabilities. Builds on `,Q(`code`,{children:`numgo`}),` for tensor operations and `,Q(`code`,{children:`gpu`}),` for hardware acceleration.`]})]}),Q(`section`,{class:`section`,id:`pgm-packages`,children:[Q(`h2`,{children:`PGM Packages (lib/pgm/)`}),Q(`p`,{children:`The PGM module provides a complete probabilistic graphical models toolkit, inspired by pgmpy. It is one of several domain modules in the datascience library.`}),Q(`h3`,{children:`base -- Foundational Graph Types`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/base`})]}),Q(`p`,{children:`Provides the underlying graph structures that all model types are built on.`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DAG`})}),Q(`td`,{children:`Directed acyclic graph with cycle detection, topological sort, d-separation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`PDAG`})}),Q(`td`,{children:`Partially directed acyclic graph for Markov equivalence classes`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`UndirectedGraph`})}),Q(`td`,{children:`Undirected graph for Markov networks`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ADMG`})}),Q(`td`,{children:`Acyclic directed mixed graph (with bidirected edges for latent confounders)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MAG`})}),Q(`td`,{children:`Maximal ancestral graph`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SimpleCausalModel`})}),Q(`td`,{children:`Basic causal model with intervention semantics`})]})]})]}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/pgm/base"

dag := base.NewDAG()
dag.AddNode("X")
dag.AddNode("Y")
dag.AddNode("Z")
dag.AddEdge("X", "Y")
dag.AddEdge("Y", "Z")

// Check d-separation: X _||_ Z | Y?
separated := dag.DSeparation([]string{"X"}, []string{"Z"}, []string{"Y"})
fmt.Println("X _||_ Z | Y:", separated) // true`})}),Q(`h3`,{children:`models -- Probabilistic Model Types`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/models`})]}),Q(`p`,{children:`13 model types for different PGM use cases:`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`}),Q(`th`,{children:`Constructor`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BayesianNetwork`})}),Q(`td`,{children:`DAG with TabularCPDs. The primary model type for most use cases.`}),Q(`td`,{children:Q(`code`,{children:`NewBayesianNetwork()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MarkovNetwork`})}),Q(`td`,{children:`Undirected graphical model with factor potentials.`}),Q(`td`,{children:Q(`code`,{children:`NewMarkovNetwork()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DynamicBayesianNetwork`})}),Q(`td`,{children:`BN over time slices for temporal modeling.`}),Q(`td`,{children:Q(`code`,{children:`NewDynamicBayesianNetwork()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NaiveBayes`})}),Q(`td`,{children:`Naive Bayes classifier (BN with single class parent).`}),Q(`td`,{children:Q(`code`,{children:`NewNaiveBayes()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SEM`})}),Q(`td`,{children:`Structural Equation Model with linear/nonlinear equations.`}),Q(`td`,{children:Q(`code`,{children:`NewSEM()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`FactorGraph`})}),Q(`td`,{children:`Bipartite graph of variable nodes and factor nodes.`}),Q(`td`,{children:Q(`code`,{children:`NewFactorGraph()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`JunctionTree`})}),Q(`td`,{children:`Clique tree for exact inference via message passing.`}),Q(`td`,{children:Q(`code`,{children:`NewJunctionTree()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ClusterGraph`})}),Q(`td`,{children:`Generalized cluster graph (superset of junction tree).`}),Q(`td`,{children:Q(`code`,{children:`NewClusterGraph()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LinearGaussianBN`})}),Q(`td`,{children:`BN with continuous, linearly-related Gaussian variables.`}),Q(`td`,{children:Q(`code`,{children:`NewLinearGaussianBN()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`FunctionalBN`})}),Q(`td`,{children:`BN where CPDs are defined by arbitrary functions.`}),Q(`td`,{children:Q(`code`,{children:`NewFunctionalBN()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MarkovChain`})}),Q(`td`,{children:`First-order Markov chain for sequential data.`}),Q(`td`,{children:Q(`code`,{children:`NewMarkovChain()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiscreteBayesianNetwork`})}),Q(`td`,{children:`Specialized discrete-only BN with optimized operations.`}),Q(`td`,{children:Q(`code`,{children:`NewDiscreteBayesianNetwork()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiscreteMarkovNetwork`})}),Q(`td`,{children:`Specialized discrete-only Markov network.`}),Q(`td`,{children:Q(`code`,{children:`NewDiscreteMarkovNetwork()`})})]})]})]}),Q(`p`,{children:[`Key methods shared by most model types: `,Q(`code`,{children:`AddNode`}),`, `,Q(`code`,{children:`AddEdge`}),`, `,Q(`code`,{children:`Nodes()`}),`,`,Q(`code`,{children:`Edges()`}),`, `,Q(`code`,{children:`SetStates`}),`, `,Q(`code`,{children:`SetCPD`}),`, `,Q(`code`,{children:`CheckModel()`}),`,`,Q(`code`,{children:`ToMarkovFactors()`}),`, `,Q(`code`,{children:`ToJunctionTree()`}),`.`]}),Q(`h3`,{children:`factors -- Factor Representations`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/factors`})]}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`}),Q(`th`,{children:`Constructor`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiscreteFactor`})}),Q(`td`,{children:`General discrete factor with product, marginalize, reduce, normalize operations.`}),Q(`td`,{children:Q(`code`,{children:`NewDiscreteFactor()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`TabularCPD`})}),Q(`td`,{children:`Conditional probability distribution table. The standard CPD type for BayesianNetwork.`}),Q(`td`,{children:Q(`code`,{children:`NewTabularCPD()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`JointProbabilityDistribution`})}),Q(`td`,{children:`Full joint distribution over a set of variables.`}),Q(`td`,{children:Q(`code`,{children:`NewJPD()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LinearGaussianCPD`})}),Q(`td`,{children:`Linear Gaussian conditional: child = sum(beta_i * parent_i) + noise.`}),Q(`td`,{children:Q(`code`,{children:`NewLinearGaussianCPD()`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NoisyOR`})}),Q(`td`,{children:`Noisy-OR parameterization. Compact CPD for nodes with many binary parents.`}),Q(`td`,{children:Q(`code`,{children:`NewNoisyOR()`})})]})]})]}),Q(`h3`,{children:`inference -- Inference Algorithms`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/inference`})]}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`}),Q(`th`,{children:`Constructor`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`VariableElimination`})}),Q(`td`,{children:[`Exact inference via factor elimination. Supports `,Q(`code`,{children:`Query()`}),` for posteriors and `,Q(`code`,{children:`MAP()`}),` for most-probable assignments.`]}),Q(`td`,{children:Q(`code`,{children:`NewVariableElimination(factors)`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BeliefPropagation`})}),Q(`td`,{children:`Message-passing on junction trees. Calibrate once, query multiple times.`}),Q(`td`,{children:Q(`code`,{children:`NewBeliefPropagation(cliques, separators, factors)`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MPLP`})}),Q(`td`,{children:`Max-Product Linear Programming for MAP inference.`}),Q(`td`,{children:Q(`code`,{children:`NewMPLP(factors)`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ApproxInference`})}),Q(`td`,{children:`Sampling-based approximate inference. Uses likelihood-weighted sampling.`}),Q(`td`,{children:Q(`code`,{children:`NewApproxInference(bn, nSamples)`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`CausalInference`})}),Q(`td`,{children:`Do-calculus interventional queries. Computes P(Y | do(X=x)).`}),Q(`td`,{children:Q(`code`,{children:`NewCausalInference(bn)`})})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DBNInference`})}),Q(`td`,{children:`Inference over dynamic Bayesian networks across time slices.`}),Q(`td`,{children:Q(`code`,{children:`NewDBNInference(dbn)`})})]})]})]}),Q(`h3`,{children:`sampling -- Sampling Methods`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/sampling`})]}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Methods`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BayesianModelSampling`})}),Q(`td`,{children:[Q(`code`,{children:`ForwardSample`}),`, `,Q(`code`,{children:`RejectionSample`}),`, `,Q(`code`,{children:`LikelihoodWeightedSample`})]}),Q(`td`,{children:`Exact and weighted sampling from BN joint distribution.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`GibbsSampling`})}),Q(`td`,{children:Q(`code`,{children:`Sample`})}),Q(`td`,{children:`MCMC Gibbs sampler with configurable burn-in and thinning.`})]})]})]}),Q(`h3`,{children:`learning -- Learning Algorithms`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/learning`})]}),Q(`p`,{children:[`15+ algorithms for parameter estimation and structure learning. See the `,Q(T,{to:`/api`,children:`API Reference`}),` for full details.`]}),Q(`h3`,{children:`Other PGM Packages`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Package`}),Q(`th`,{children:`Import Path`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ci_tests`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/ci_tests`})}),Q(`td`,{children:`16 conditional independence tests (ChiSquare, FisherZ, GCM, etc.)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`structure_score`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/structure_score`})}),Q(`td`,{children:`13 scoring functions (BIC, AIC, BDeu, K2, etc.)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`identification`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/identification`})}),Q(`td`,{children:`Causal effect identification (back-door, front-door)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`prediction`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/prediction`})}),Q(`td`,{children:`DoubleML, naive adjustment, IV regression`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`metrics`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/metrics`})}),Q(`td`,{children:`SHD, confusion matrices, Fisher's C, log-likelihood`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`independencies`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/independencies`})}),Q(`td`,{children:`Independence assertion representations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`readwrite`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/readwrite`})}),Q(`td`,{children:`10 file format readers and writers`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`config`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/config`})}),Q(`td`,{children:`Global PGM configuration`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`utils`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/utils`})}),Q(`td`,{children:`Shared utilities`})]})]})]})]}),Q(`section`,{class:`section`,id:`file-formats`,children:[Q(`h2`,{children:`File Formats`}),Q(`p`,{children:[`datascience supports 10 file formats for reading and writing probabilistic graphical models. All formats are accessed through the `,Q(`code`,{children:`readwrite`}),` package.`]}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Format`}),Q(`th`,{children:`Extension`}),Q(`th`,{children:`Read`}),Q(`th`,{children:`Write`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`BIF`})}),Q(`td`,{children:`.bif`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Bayesian Interchange Format. The standard PGM format. Stores structure, states, and CPD tables in a human-readable text format.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`XMLBIF`})}),Q(`td`,{children:`.xmlbif`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`XML-based BIF. Same information as BIF but in XML for easier parsing by other tools.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`NET`})}),Q(`td`,{children:`.net`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Hugin NET format. Used by the Hugin BN software.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`UAI`})}),Q(`td`,{children:`.uai`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`UAI format. Used by the UAI inference competition. Compact numeric format.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`XDSL`})}),Q(`td`,{children:`.xdsl`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`GeNIe XDSL format. Used by the GeNIe/SMILE BN software.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`PomdpX`})}),Q(`td`,{children:`.pomdpx`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`--`}),Q(`td`,{children:`POMDP XML format. Read-only. Used in POMDP planning literature.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`XBN`})}),Q(`td`,{children:`.xbn`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`--`}),Q(`td`,{children:`Microsoft XBN format. Read-only. Legacy Microsoft Research format.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`CSV`})}),Q(`td`,{children:`.csv`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`CSV model serialization. Stores structure and parameters in CSV tables.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`JSON`})}),Q(`td`,{children:`.json`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`JSON model serialization. Ideal for web applications and REST APIs.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`XML`})}),Q(`td`,{children:`.xml`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`XML model serialization. General-purpose XML format.`})]})]})]}),Q(`h3`,{children:`BIF Format Example`}),Q(`pre`,{children:Q(`code`,{children:`network asia {
}

variable VisitAsia {
  type discrete [ 2 ] { no, yes };
}

variable Tuberculosis {
  type discrete [ 2 ] { no, yes };
}

probability ( VisitAsia ) {
  table 0.99, 0.01;
}

probability ( Tuberculosis | VisitAsia ) {
  (no) 0.99, 0.01;
  (yes) 0.95, 0.05;
}`})}),Q(`h3`,{children:`JSON Format Example`}),Q(`pre`,{children:Q(`code`,{children:`{
  "nodes": ["A", "B", "C"],
  "edges": [["A", "B"], ["B", "C"]],
  "states": {
    "A": ["a0", "a1"],
    "B": ["b0", "b1"],
    "C": ["c0", "c1"]
  },
  "cpds": {
    "A": {
      "variable": "A",
      "cardinality": 2,
      "values": [0.4, 0.6],
      "parents": [],
      "parent_cardinalities": []
    }
  }
}`})}),Q(`h3`,{children:`Reading and Writing`}),Q(`pre`,{children:Q(`code`,{children:`import (
    "os"
    "github.com/asymmetric-effort/datascience/lib/pgm/readwrite"
)

// Read BIF
f, _ := os.Open("model.bif")
bn, _ := readwrite.ReadBIF(f)
f.Close()

// Write as JSON (for a web API)
out, _ := os.Create("model.json")
readwrite.WriteJSONModel(out, bn)
out.Close()

// Convert: read NET, write XMLBIF
netFile, _ := os.Open("model.net")
bn2, _ := readwrite.ReadNET(netFile)
netFile.Close()

xmlbifFile, _ := os.Create("model.xmlbif")
readwrite.WriteXMLBIF(xmlbifFile, bn2)
xmlbifFile.Close()`})})]}),Q(`section`,{class:`section`,id:`example-models`,children:[Q(`h2`,{children:`Example Models`}),Q(`p`,{children:[`datascience ships with 25 built-in Bayesian networks accessible via the `,Q(`code`,{children:`example_models`}),` package. 13 models include full CPDs (conditional probability distributions); 12 are structure-only for use in learning and benchmarking.`]}),Q(`h3`,{children:`Models with Full CPDs`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Name`}),Q(`th`,{children:`Description`}),Q(`th`,{children:`Nodes`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`student`})}),Q(`td`,{children:`Classic Student network (difficulty, intelligence, grade, SAT, letter)`}),Q(`td`,{children:`5`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`asia`})}),Q(`td`,{children:`Lung disease diagnosis`}),Q(`td`,{children:`8`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`alarm`})}),Q(`td`,{children:`Monitoring system with alarm, burglary, earthquake`}),Q(`td`,{children:`5`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`cancer`})}),Q(`td`,{children:`Cancer diagnosis network`}),Q(`td`,{children:`5`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`watersprinkler`})}),Q(`td`,{children:`Classic sprinkler/rain/wet grass example`}),Q(`td`,{children:`4`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`survey`})}),Q(`td`,{children:`Survey response model`}),Q(`td`,{children:`6`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`montyhall`})}),Q(`td`,{children:`Monty Hall problem as a Bayesian network`}),Q(`td`,{children:`3`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`dogproblem`})}),Q(`td`,{children:`Dog behavior inference`}),Q(`td`,{children:`5`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`frauddetection`})}),Q(`td`,{children:`Financial fraud detection model`}),Q(`td`,{children:`5`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`medicaldiagnosis`})}),Q(`td`,{children:`Medical symptom/disease model`}),Q(`td`,{children:`8`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`earthquake`})}),Q(`td`,{children:`Earthquake alert network`}),Q(`td`,{children:`5`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`visitasia`})}),Q(`td`,{children:`Visit to Asia variant`}),Q(`td`,{children:`8`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`cointoss`})}),Q(`td`,{children:`Simple coin toss model`}),Q(`td`,{children:`2`})]})]})]}),Q(`h3`,{children:`Structure-Only Models (Large Networks)`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Name`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`sachs`})}),Q(`td`,{children:`Protein signaling network (11 nodes)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`child`})}),Q(`td`,{children:`Child health assessment network`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`insurance`})}),Q(`td`,{children:`Insurance risk assessment`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`alarmfull`})}),Q(`td`,{children:`Full ALARM monitoring network (37 nodes)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`water`})}),Q(`td`,{children:`Water treatment network`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`mildew`})}),Q(`td`,{children:`Crop disease model`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`barley`})}),Q(`td`,{children:`Barley crop yield model`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`hailfinder`})}),Q(`td`,{children:`Severe weather prediction`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`hepar2`})}),Q(`td`,{children:`Liver disorder diagnosis`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`win95pts`})}),Q(`td`,{children:`Windows 95 printer troubleshooting`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`pathfinder`})}),Q(`td`,{children:`Pathology diagnosis`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`pigs`})}),Q(`td`,{children:`Pig breeding network`})]})]})]}),Q(`h3`,{children:`Usage`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/example_models"

// List all available models
names := example_models.List()
for _, name := range names {
    fmt.Println(name)
}

// Load a specific model
bn, err := example_models.Get("asia")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Nodes: %d, Edges: %d\\n", len(bn.Nodes()), len(bn.Edges()))

// Use the model for inference
facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)
result, _ := ve.Query([]string{"Dyspnea"}, map[string]int{"Smoker": 1})
fmt.Println(result.Values().Data())`})})]}),Q(`section`,{class:`section`,id:`datasets`,children:[Q(`h2`,{children:`Datasets`}),Q(`p`,{children:[`datascience includes 40 built-in datasets accessible via the `,Q(`code`,{children:`examples/datasets`}),` package. These datasets are embedded in the binary using Go's `,Q(`code`,{children:`embed`}),` package, so they are always available without external file dependencies.`]}),Q(`h3`,{children:`BN-Specific Datasets`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Name`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`asia`})}),Q(`td`,{children:`Sampled data from the Asia (lung disease) network`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`alarm`})}),Q(`td`,{children:`Sampled data from the ALARM monitoring network`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`sachs`})}),Q(`td`,{children:`Protein signaling data (Sachs et al.)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`cancer`})}),Q(`td`,{children:`Cancer diagnosis observations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`student`})}),Q(`td`,{children:`Student performance data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`sprinkler`})}),Q(`td`,{children:`Sprinkler/rain/wet grass observations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`survey`})}),Q(`td`,{children:`Survey response data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`earthquake`})}),Q(`td`,{children:`Earthquake alert observations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`child`})}),Q(`td`,{children:`Child health assessment data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`insurance`})}),Q(`td`,{children:`Insurance risk data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`water`})}),Q(`td`,{children:`Water treatment data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`mildew`})}),Q(`td`,{children:`Crop disease data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`hailfinder`})}),Q(`td`,{children:`Severe weather observations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`hepar2`})}),Q(`td`,{children:`Liver disorder data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`barley`})}),Q(`td`,{children:`Barley crop data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`win95pts`})}),Q(`td`,{children:`Windows 95 troubleshooting data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`andes`})}),Q(`td`,{children:`ANDES intelligent tutoring system data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`munin`})}),Q(`td`,{children:`MUNIN neural network data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lucas`})}),Q(`td`,{children:`LUCAS causal discovery benchmark`})]})]})]}),Q(`h3`,{children:`Classic ML Datasets`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Name`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`titanic`})}),Q(`td`,{children:`Titanic survival data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`iris`})}),Q(`td`,{children:`Fisher's Iris flower dataset`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`heart`})}),Q(`td`,{children:`Heart disease prediction`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`wine`})}),Q(`td`,{children:`Wine quality classification`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`boston`})}),Q(`td`,{children:`Boston housing prices`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`pima_diabetes`})}),Q(`td`,{children:`Pima Indians diabetes`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`adult`})}),Q(`td`,{children:`Adult income prediction (Census)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`breast_cancer`})}),Q(`td`,{children:`Wisconsin breast cancer diagnosis`})]})]})]}),Q(`h3`,{children:`UCI Repository Datasets`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Name`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`zoo`})}),Q(`td`,{children:`Zoo animal classification`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`glass`})}),Q(`td`,{children:`Glass identification`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ecoli`})}),Q(`td`,{children:`E. coli protein localization`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`monks`})}),Q(`td`,{children:`MONKS problem`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`nursery`})}),Q(`td`,{children:`Nursery school evaluation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`credit_approval`})}),Q(`td`,{children:`Credit card approval`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`balance_scale`})}),Q(`td`,{children:`Balance scale weight/distance`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`automobile`})}),Q(`td`,{children:`Automobile price prediction`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`mushroom`})}),Q(`td`,{children:`Mushroom edibility classification`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`car_evaluation`})}),Q(`td`,{children:`Car evaluation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`hepatitis`})}),Q(`td`,{children:`Hepatitis prognosis`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`vote`})}),Q(`td`,{children:`Congressional voting records`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`tic_tac_toe`})}),Q(`td`,{children:`Tic-tac-toe endgame`})]})]})]}),Q(`h3`,{children:`Usage`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/examples/datasets"

// List all available datasets
names := datasets.List()

// Load a dataset as a DataFrame
df, err := datasets.Load("asia")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Rows: %d, Columns: %d\\n", df.NRows(), len(df.Columns()))

// Use the dataset for structure learning
scorer := structure_score.NewBIC()
hc := learning.NewHillClimbSearch(df, scorer.LocalScore)
bn, _ := hc.Estimate()`})})]}),Q(`section`,{class:`section`,id:`testing`,children:[Q(`h2`,{children:`Testing`}),Q(`p`,{children:`datascience has approximately 5,000 tests and 392 cross-validation fixtures spanning unit tests, integration tests, and cross-validation tests across the library.`}),Q(`h3`,{children:`Running Tests`}),Q(`pre`,{children:Q(`code`,{children:`# Run all tests
go test ./...

# Run tests for a specific package
go test ./lib/pgm/inference/...

# Run with verbose output
go test -v ./lib/pgm/models/...

# Run with race detector
go test -race ./...

# Run a specific test function
go test -run TestVariableElimination ./lib/pgm/inference/...

# Run cross-validation tests only
go test -run CrossVal ./...

# Run benchmarks
go test -bench=. ./lib/pgm/inference/...`})}),Q(`h3`,{children:`Cross-Validation System`}),Q(`p`,{children:[`Many packages include `,Q(`code`,{children:`crossval_*_test.go`}),` files that validate algorithms against known-correct results. These tests load built-in example models, run computations, and compare outputs against pre-computed reference values. Examples:`]}),Q(`ul`,{children:[Q(`li`,{children:[Q(`code`,{children:`lib/pgm/models/crossval_dsep_test.go`}),` -- validates d-separation queries`]}),Q(`li`,{children:[Q(`code`,{children:`lib/pgm/inference/crossval_causal_test.go`}),` -- validates causal inference results`]}),Q(`li`,{children:[Q(`code`,{children:`lib/pgm/inference/crossval_ve_test.go`}),` -- validates Variable Elimination posteriors`]}),Q(`li`,{children:[Q(`code`,{children:`lib/pgm/learning/crossval_hillclimb_test.go`}),` -- validates structure learning output`]}),Q(`li`,{children:[Q(`code`,{children:`lib/pgm/sampling/crossval_forward_test.go`}),` -- validates sampling distributions`]})]}),Q(`p`,{children:`Cross-validation fixtures are generated by running the equivalent pgmpy code in Python and storing the results. This ensures datascience produces the same numerical outputs as the reference implementation.`}),Q(`h3`,{children:`Test Fixture Generation`}),Q(`p`,{children:[`Test fixtures use the built-in example models from the `,Q(`code`,{children:`example_models`}),` package. This ensures reproducible test data without external file dependencies. When adding new tests, use existing models or create new ones in the `,Q(`code`,{children:`example_models`}),` package.`]}),Q(`h3`,{children:`Writing Tests`}),Q(`pre`,{children:Q(`code`,{children:`func TestMyFeature(t *testing.T) {
    // Load a known model
    bn, err := example_models.Get("asia")
    if err != nil {
        t.Fatal(err)
    }

    // Perform computation
    facs, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    result, err := ve.Query([]string{"Dyspnea"}, map[string]int{"Smoker": 1})
    if err != nil {
        t.Fatal(err)
    }

    // Compare against expected values
    values := result.Values().Data()
    if math.Abs(values[0] - 0.304) > 0.01 {
        t.Errorf("expected P(Dyspnea=0|Smoker=1) ~ 0.304, got %f", values[0])
    }
}`})})]}),Q(`section`,{class:`section`,id:`configuration`,children:[Q(`h2`,{children:`Configuration`}),Q(`p`,{children:[`The `,Q(`code`,{children:`config`}),` package provides global configuration options that control default behavior across the PGM module.`]}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/pgm/config"

// Get the global config
cfg := config.Global()

// Configuration is used internally by various packages
// to control default inference methods, scoring functions,
// numerical tolerances, and other global settings.`})}),Q(`p`,{children:`Configuration options include numerical tolerances for probability comparisons, default inference methods, default scoring functions for structure learning, convergence thresholds for iterative algorithms (EM, Belief Propagation), and logging verbosity.`})]}),Q(`section`,{class:`section`,id:`contributing`,children:[Q(`h2`,{children:`Contributing`}),Q(`p`,{children:[`Contributions are welcome. See the full`,` `,Q(`a`,{href:`https://github.com/asymmetric-effort/datascience/blob/main/CONTRIBUTING.md`,target:`_blank`,rel:`noopener noreferrer`,children:`CONTRIBUTING.md`}),` `,`for details.`]}),Q(`h3`,{children:`Development Workflow`}),Q(`ol`,{children:[Q(`li`,{children:`Fork the repository and clone your fork`}),Q(`li`,{children:[`Create a feature branch: `,Q(`code`,{children:`git checkout -b feature/my-feature`})]}),Q(`li`,{children:`Make changes and add tests`}),Q(`li`,{children:[`Run `,Q(`code`,{children:`go test ./...`}),` to verify all tests pass`]}),Q(`li`,{children:[`Run `,Q(`code`,{children:`go vet ./...`}),` for static analysis`]}),Q(`li`,{children:`Commit with a clear message and submit a pull request`})]}),Q(`h3`,{children:`Guidelines`}),Q(`ul`,{children:[Q(`li`,{children:[Q(`strong`,{children:`Near-zero dependencies:`}),` Do not add third-party modules. All functionality must be implemented in pure Go using only the standard library.`]}),Q(`li`,{children:[Q(`strong`,{children:`Tests required:`}),` All new functionality must include unit tests. Cross-validation tests against pgmpy are strongly encouraged.`]}),Q(`li`,{children:[Q(`strong`,{children:`Backward compatibility:`}),` Public API changes require discussion in an issue before implementation.`]}),Q(`li`,{children:[Q(`strong`,{children:`Documentation:`}),` Exported types and functions must have godoc comments.`]})]})]})]})}function Xn(){return y({title:`Libraries — datascience`,description:`Overview of all library packages in datascience: numgo, scigo, graphgo, tabgo, gpu, pgm, and tensorflow.`,canonical:`https://datascience.asymmetric-effort.com/#/libraries`}),Q(`div`,{class:`page`,children:[Q(`h1`,{children:`Libraries`}),Q(`nav`,{class:`page-toc`,children:[Q(`strong`,{children:`Packages:`}),` `,Q($,{to:`lib-overview`,children:`Overview`}),` | `,Q($,{to:`lib-numgo`,children:`numgo`}),` | `,Q($,{to:`lib-scigo`,children:`scigo`}),` | `,Q($,{to:`lib-graphgo`,children:`graphgo`}),` | `,Q($,{to:`lib-tabgo`,children:`tabgo`}),` | `,Q($,{to:`lib-gpu`,children:`gpu`}),` | `,Q($,{to:`lib-pgm`,children:`pgm`}),` | `,Q($,{to:`lib-tensorflow`,children:`tensorflow`})]}),Q(`section`,{class:`section`,id:`lib-overview`,children:[Q(`h2`,{children:`Overview`}),Q(`p`,{children:[`The datascience library is organized into independent, composable packages under `,Q(`code`,{children:`lib/`}),`. Each package can be imported and used on its own. The primitive packages (numgo, scigo, graphgo, tabgo) provide general-purpose data science foundations, while the domain packages (pgm, tensorflow) build on them for specific use cases.`]}),Q(`pre`,{children:Q(`code`,{children:`go get github.com/asymmetric-effort/datascience`})}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Package`}),Q(`th`,{children:`Import Path`}),Q(`th`,{children:`Description`}),Q(`th`,{children:`Python Equivalent`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`numgo`})}),Q(`td`,{children:Q(`code`,{children:`lib/numgo`})}),Q(`td`,{children:`N-dimensional arrays, matrices, vectors, BLAS`}),Q(`td`,{children:`numpy`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`scigo`})}),Q(`td`,{children:Q(`code`,{children:`lib/scigo`})}),Q(`td`,{children:`Distributions, optimization, special functions`}),Q(`td`,{children:`scipy`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`graphgo`})}),Q(`td`,{children:Q(`code`,{children:`lib/graphgo`})}),Q(`td`,{children:`Directed/undirected graphs, algorithms`}),Q(`td`,{children:`networkx`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`tabgo`})}),Q(`td`,{children:Q(`code`,{children:`lib/tabgo`})}),Q(`td`,{children:`DataFrames, Series, CSV/Parquet I/O`}),Q(`td`,{children:`pandas`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`gpu`})}),Q(`td`,{children:Q(`code`,{children:`lib/gpu`})}),Q(`td`,{children:`GPU compute backend`}),Q(`td`,{children:`cupy / CUDA`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`pgm`})}),Q(`td`,{children:Q(`code`,{children:`lib/pgm/...`})}),Q(`td`,{children:`Probabilistic graphical models`}),Q(`td`,{children:`pgmpy`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`strong`,{children:`tensorflow`})}),Q(`td`,{children:Q(`code`,{children:`lib/tensorflow`})}),Q(`td`,{children:`Deep learning`}),Q(`td`,{children:`tensorflow`})]})]})]})]}),Q(`section`,{class:`section`,id:`lib-numgo`,children:[Q(`h2`,{children:`numgo -- Arrays, Matrices, and BLAS`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/numgo`})]}),Q(`p`,{children:`Pure Go implementation of N-dimensional arrays, 2D matrices, and 1D vectors with BLAS-level linear algebra operations. Provides the numerical foundation for all other packages.`}),Q(`h3`,{children:`Key Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NDArray`})}),Q(`td`,{children:`N-dimensional array with shape, stride, broadcasting, element-wise operations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Matrix`})}),Q(`td`,{children:`2D matrix with multiply, transpose, inverse, determinant, eigenvalues, decompositions`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Vector`})}),Q(`td`,{children:`1D vector with dot product, norm, cross product, element-wise operations`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/numgo"

// Matrix operations
m := numgo.NewMatrixFromData(3, 3, []float64{
    1, 2, 3,
    0, 1, 4,
    5, 6, 0,
})
det := m.Det()           // determinant
inv := m.Inverse()       // matrix inverse
t := m.Transpose()       // transpose

// Vector operations
v1 := numgo.NewVector([]float64{1, 2, 3})
v2 := numgo.NewVector([]float64{4, 5, 6})
dot := v1.Dot(v2)       // 32
norm := v1.Norm()        // sqrt(14)

// NDArray operations
arr := numgo.NewNDArray([]int{2, 3, 4})
arr.Fill(1.0)
sum := arr.Sum()         // 24`})})]}),Q(`section`,{class:`section`,id:`lib-scigo`,children:[Q(`h2`,{children:`scigo -- Statistical Computing`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/scigo`})]}),Q(`p`,{children:`Statistical distributions (Normal, Chi-Squared, Beta, Gamma, Student-t, Uniform, Exponential), optimization routines (gradient descent, Newton's method, minimization), special functions (gamma, beta, digamma), and hypothesis tests.`}),Q(`h3`,{children:`Key Capabilities`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Category`}),Q(`th`,{children:`Functions / Types`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`Distributions`}),Q(`td`,{children:[Q(`code`,{children:`Normal`}),`, `,Q(`code`,{children:`ChiSquared`}),`, `,Q(`code`,{children:`Beta`}),`, `,Q(`code`,{children:`Gamma`}),`, `,Q(`code`,{children:`StudentT`}),`, `,Q(`code`,{children:`Uniform`}),`, `,Q(`code`,{children:`Exponential`}),` -- each with PDF, CDF, PPF, Sample`]})]}),Q(`tr`,{children:[Q(`td`,{children:`Optimization`}),Q(`td`,{children:[Q(`code`,{children:`Minimize`}),`, `,Q(`code`,{children:`GradientDescent`}),`, `,Q(`code`,{children:`NewtonMethod`})]})]}),Q(`tr`,{children:[Q(`td`,{children:`Statistics`}),Q(`td`,{children:[Q(`code`,{children:`Mean`}),`, `,Q(`code`,{children:`Std`}),`, `,Q(`code`,{children:`PearsonCorrelation`}),`, p-value computation`]})]}),Q(`tr`,{children:[Q(`td`,{children:`Special Functions`}),Q(`td`,{children:`Gamma, beta, incomplete gamma, digamma, log-gamma`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/scigo"

// Normal distribution
n := scigo.NewNormal(0, 1)
cdf := n.CDF(1.96)       // ~0.975
ppf := n.PPF(0.975)      // ~1.96

// Chi-squared test
chi2 := scigo.NewChiSquared(5)
pValue := 1.0 - chi2.CDF(11.07)

// Optimization
result := scigo.Minimize(func(x float64) float64 {
    return (x - 3) * (x - 3)
}, 0.0, 10.0)  // ~3.0`})})]}),Q(`section`,{class:`section`,id:`lib-graphgo`,children:[Q(`h2`,{children:`graphgo -- Graph Algorithms`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/graphgo`})]}),Q(`p`,{children:`Directed and undirected graphs with a full suite of algorithms: topological sort, d-separation, moral graph, triangulation, maximum cardinality search, clique finding, connected components, shortest paths, and cycle detection.`}),Q(`h3`,{children:`Key Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiGraph`})}),Q(`td`,{children:`Directed graph with successors, predecessors, topological sort, d-separation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Graph`})}),Q(`td`,{children:`Undirected graph with neighbors, degree, connected components`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`PDAG`})}),Q(`td`,{children:`Partially directed acyclic graph for equivalence classes`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/graphgo"

g := graphgo.NewDiGraph()
g.AddNode("A")
g.AddNode("B")
g.AddNode("C")
g.AddEdge("A", "B")
g.AddEdge("B", "C")

sorted := g.TopologicalSort()  // [A, B, C]
hasCycle := g.HasCycle()       // false`})})]}),Q(`section`,{class:`section`,id:`lib-tabgo`,children:[Q(`h2`,{children:`tabgo -- Tabular Data`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/tabgo`})]}),Q(`p`,{children:`DataFrames and Series for tabular data manipulation. Supports CSV, Parquet, and Excel I/O, row filtering, groupby, merge, column operations, and statistical summaries.`}),Q(`h3`,{children:`Key Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DataFrame`})}),Q(`td`,{children:`Tabular data with named columns, filtering, groupby, merge`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Series`})}),Q(`td`,{children:`Single column with value counts, unique values, statistics`})]})]})]}),Q(`h3`,{children:`I/O Functions`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Function`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadCSV / WriteCSV`})}),Q(`td`,{children:`CSV file I/O`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadParquet / WriteParquet`})}),Q(`td`,{children:`Apache Parquet I/O`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadXLSX`})}),Q(`td`,{children:`Excel file reading`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/tabgo"

df, _ := tabgo.ReadCSV("data.csv")
fmt.Printf("%d rows, %d columns\\n", df.NRows(), len(df.Columns()))

col := df.Column("Temperature")
fmt.Println("Unique:", col.Unique())

filtered := df.Filter(func(row map[string]interface{}) bool {
    return row["Temperature"].(int) > 70
})
tabgo.WriteCSV(filtered, "warm_days.csv")`})})]}),Q(`section`,{class:`section`,id:`lib-gpu`,children:[Q(`h2`,{children:`gpu -- GPU Compute Backend`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/gpu`})]}),Q(`p`,{children:`Optional GPU acceleration for compute-intensive operations. Provides GPU-backed matrix operations, tensor computations, and factor operations that can significantly speed up inference, learning, and deep learning training on large models.`}),Q(`p`,{children:`The GPU package is designed to be a drop-in accelerator. Other packages detect GPU availability and automatically offload heavy computations when a GPU is present.`})]}),Q(`section`,{class:`section`,id:`lib-pgm`,children:[Q(`h2`,{children:`pgm -- Probabilistic Graphical Models`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/...`})]}),Q(`p`,{children:[`A complete PGM toolkit inspired by `,Q(`a`,{href:`https://pgmpy.org`,target:`_blank`,rel:`noopener noreferrer`,children:`pgmpy`}),`. Provides 13 model types, 7 inference algorithms, 15+ learning algorithms, 16 conditional independence tests, 13 scoring functions, 10 file formats, causal inference, and more.`]}),Q(`h3`,{children:`Sub-packages`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Package`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/base`})}),Q(`td`,{children:`DAG, PDAG, UndirectedGraph, ADMG, MAG, SimpleCausalModel`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/models`})}),Q(`td`,{children:`13 probabilistic model types (BayesianNetwork, MarkovNetwork, etc.)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/factors`})}),Q(`td`,{children:`TabularCPD, DiscreteFactor, LinearGaussianCPD, NoisyOR, JPD`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/inference`})}),Q(`td`,{children:`Variable Elimination, Belief Propagation, MPLP, Causal, Approximate, DBN`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/sampling`})}),Q(`td`,{children:`Forward, rejection, likelihood-weighted, Gibbs sampling`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/learning`})}),Q(`td`,{children:`MLE, Bayesian, EM, HillClimb, PC, GES, ExhaustiveSearch, TreeSearch, MMHC`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/ci_tests`})}),Q(`td`,{children:`ChiSquare, FisherZ, GCM, and 13 more CI tests`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/structure_score`})}),Q(`td`,{children:`BIC, AIC, BDeu, BDs, K2, LogLikelihood, Gaussian scores`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/identification`})}),Q(`td`,{children:`Back-door and front-door causal identification`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/prediction`})}),Q(`td`,{children:`DoubleML, naive adjustment, IV regression`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/metrics`})}),Q(`td`,{children:`SHD, confusion matrices, Fisher's C, log-likelihood`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/readwrite`})}),Q(`td`,{children:`BIF, XMLBIF, NET, UAI, XDSL, JSON, CSV, XML I/O`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/independencies`})}),Q(`td`,{children:`Independence assertion representations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/config`})}),Q(`td`,{children:`Global configuration`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`lib/pgm/utils`})}),Q(`td`,{children:`Shared utilities`})]})]})]}),Q(`h3`,{children:`Quick Example`}),Q(`pre`,{children:Q(`code`,{children:`import (
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

bn := models.NewBayesianNetwork()
bn.AddNode("A")
bn.AddNode("B")
bn.AddEdge("A", "B")
bn.SetStates("A", []string{"a0", "a1"})
bn.SetStates("B", []string{"b0", "b1"})
bn.SetCPD("A", factors.NewTabularCPD("A", 2, []float64{0.6, 0.4}, nil, nil))
bn.SetCPD("B", factors.NewTabularCPD("B", 2, []float64{0.2, 0.8, 0.75, 0.25}, []string{"A"}, []int{2}))

facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)
result, _ := ve.Query([]string{"B"}, map[string]int{"A": 1})
fmt.Println("P(B | A=1):", result.Values().Data())`})}),Q(`p`,{children:[`See the `,Q(T,{to:`/docs`,children:`Documentation`}),` and `,Q(T,{to:`/api`,children:`API Reference`}),` for complete details.`]})]}),Q(`section`,{class:`section`,id:`lib-tensorflow`,children:[Q(`h2`,{children:`tensorflow -- Deep Learning`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/tensorflow`})]}),Q(`p`,{children:[`TensorFlow-compatible deep learning in pure Go. Provides neural network construction, training, and inference capabilities. Builds on `,Q(`code`,{children:`numgo`}),` for tensor operations and `,Q(`code`,{children:`gpu`}),` for optional hardware acceleration.`]}),Q(`p`,{children:[`The tensorflow package brings deep learning capabilities to the datascience library, complementing the probabilistic graphical models in `,Q(`code`,{children:`lib/pgm`}),`. Together, they cover both classical statistical modeling and modern neural network approaches.`]})]})]})}function Zn(){return y({title:`API Reference — datascience`,description:`Complete Go API reference for datascience: PGM models, factors, inference, sampling, learning, TensorFlow, BLAS, financial modeling, and library packages.`,canonical:`https://datascience.asymmetric-effort.com/#/api`}),Q(`div`,{class:`page`,children:[Q(`h1`,{children:`API Reference`}),Q(`nav`,{class:`page-toc`,children:[Q(`strong`,{children:`Packages:`}),` `,Q($,{to:`import-paths`,children:`Import Paths`}),` | `,Q($,{to:`api-models`,children:`Models`}),` | `,Q($,{to:`api-factors`,children:`Factors`}),` | `,Q($,{to:`api-inference`,children:`Inference`}),` | `,Q($,{to:`api-sampling`,children:`Sampling`}),` | `,Q($,{to:`api-learning`,children:`Learning`}),` | `,Q($,{to:`api-ci-tests`,children:`CI Tests`}),` | `,Q($,{to:`api-structure-score`,children:`Structure Scores`}),` | `,Q($,{to:`api-identification`,children:`Identification`}),` | `,Q($,{to:`api-prediction`,children:`Prediction`}),` | `,Q($,{to:`api-metrics`,children:`Metrics`}),` | `,Q($,{to:`api-readwrite`,children:`Readwrite`}),` | `,Q($,{to:`api-base`,children:`Base`}),` | `,Q($,{to:`api-independencies`,children:`Independencies`}),` | `,Q($,{to:`api-config`,children:`Config`}),` | `,Q($,{to:`api-utils`,children:`Utils`}),` | `,Q($,{to:`api-numgo`,children:`numgo`}),` | `,Q($,{to:`api-scigo`,children:`scigo`}),` | `,Q($,{to:`api-graphgo`,children:`graphgo`}),` | `,Q($,{to:`api-tabgo`,children:`tabgo`}),` | `,Q($,{to:`api-gpu`,children:`gpu`}),` | `,Q($,{to:`api-tensorflow`,children:`tensorflow`}),` | `,Q($,{to:`api-example-models`,children:`example_models`})]}),Q(`section`,{class:`section`,id:`import-paths`,children:[Q(`h2`,{children:`Import Paths`}),Q(`pre`,{children:Q(`code`,{children:`import (
    // PGM packages (lib/pgm/)
    "github.com/asymmetric-effort/datascience/lib/pgm/base"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
    "github.com/asymmetric-effort/datascience/lib/pgm/sampling"
    "github.com/asymmetric-effort/datascience/lib/pgm/learning"
    "github.com/asymmetric-effort/datascience/lib/pgm/ci_tests"
    "github.com/asymmetric-effort/datascience/lib/pgm/structure_score"
    "github.com/asymmetric-effort/datascience/lib/pgm/identification"
    "github.com/asymmetric-effort/datascience/lib/pgm/prediction"
    "github.com/asymmetric-effort/datascience/lib/pgm/metrics"
    "github.com/asymmetric-effort/datascience/lib/pgm/independencies"
    "github.com/asymmetric-effort/datascience/lib/pgm/readwrite"
    "github.com/asymmetric-effort/datascience/lib/pgm/config"
    "github.com/asymmetric-effort/datascience/lib/pgm/utils"

    // Primitive library packages (lib/)
    "github.com/asymmetric-effort/datascience/lib/numgo"
    "github.com/asymmetric-effort/datascience/lib/scigo"
    "github.com/asymmetric-effort/datascience/lib/graphgo"
    "github.com/asymmetric-effort/datascience/lib/tabgo"
    "github.com/asymmetric-effort/datascience/lib/gpu"

    // Deep learning (lib/tensorflow/)
    "github.com/asymmetric-effort/datascience/lib/tensorflow"

    // Built-in example models and datasets
    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/examples/datasets"
)`})})]}),Q(`section`,{class:`section`,id:`api-models`,children:[Q(`h2`,{children:`models -- Probabilistic Model Types`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/models`})]}),Q(`p`,{children:`Provides 13 probabilistic graphical model types. All model types share a common interface for node/edge management, state definition, and CPD assignment.`}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BayesianNetwork`})}),Q(`td`,{children:`Directed acyclic graphical model with TabularCPDs. The primary model type.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MarkovNetwork`})}),Q(`td`,{children:`Undirected graphical model with factor potentials.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DynamicBayesianNetwork`})}),Q(`td`,{children:`BN over time slices for temporal modeling.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NaiveBayes`})}),Q(`td`,{children:`Naive Bayes classifier (BN with single class parent).`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SEM`})}),Q(`td`,{children:`Structural Equation Model with linear/nonlinear equations.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`FactorGraph`})}),Q(`td`,{children:`Bipartite graph of variable nodes and factor nodes.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`JunctionTree`})}),Q(`td`,{children:`Clique tree for exact inference via message passing.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ClusterGraph`})}),Q(`td`,{children:`Generalized cluster graph (superset of junction tree).`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LinearGaussianBN`})}),Q(`td`,{children:`BN with continuous linearly-related Gaussian variables.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`FunctionalBN`})}),Q(`td`,{children:`BN where CPDs are defined by arbitrary Go functions.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MarkovChain`})}),Q(`td`,{children:`First-order Markov chain for sequential data.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiscreteBayesianNetwork`})}),Q(`td`,{children:`Specialized discrete-only BN with optimized operations.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiscreteMarkovNetwork`})}),Q(`td`,{children:`Specialized discrete-only Markov network.`})]})]})]}),Q(`h3`,{children:`Constructors`}),Q(`pre`,{children:Q(`code`,{children:`bn := models.NewBayesianNetwork()
mn := models.NewMarkovNetwork()
dbn := models.NewDynamicBayesianNetwork()
nb := models.NewNaiveBayes()
sem := models.NewSEM()
fg := models.NewFactorGraph()
jt := models.NewJunctionTree()
cg := models.NewClusterGraph()
lgbn := models.NewLinearGaussianBN()
fbn := models.NewFunctionalBN()
mc := models.NewMarkovChain()
dbn2 := models.NewDiscreteBayesianNetwork()
dmn := models.NewDiscreteMarkovNetwork()`})}),Q(`h3`,{children:`BayesianNetwork Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`AddNode`})}),Q(`td`,{children:Q(`code`,{children:`(name string)`})}),Q(`td`,{children:`Add a node to the network`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`AddEdge`})}),Q(`td`,{children:Q(`code`,{children:`(from, to string)`})}),Q(`td`,{children:`Add a directed edge`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`RemoveNode`})}),Q(`td`,{children:Q(`code`,{children:`(name string)`})}),Q(`td`,{children:`Remove a node and its edges`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`RemoveEdge`})}),Q(`td`,{children:Q(`code`,{children:`(from, to string)`})}),Q(`td`,{children:`Remove a directed edge`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Nodes`})}),Q(`td`,{children:Q(`code`,{children:`() []string`})}),Q(`td`,{children:`List all node names`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Edges`})}),Q(`td`,{children:Q(`code`,{children:`() [][2]string`})}),Q(`td`,{children:`List all edges as [from, to] pairs`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Parents`})}),Q(`td`,{children:Q(`code`,{children:`(name string) []string`})}),Q(`td`,{children:`Get parents of a node`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Children`})}),Q(`td`,{children:Q(`code`,{children:`(name string) []string`})}),Q(`td`,{children:`Get children of a node`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SetStates`})}),Q(`td`,{children:Q(`code`,{children:`(name string, states []string)`})}),Q(`td`,{children:`Define possible values for a node`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`GetStates`})}),Q(`td`,{children:Q(`code`,{children:`(name string) []string`})}),Q(`td`,{children:`Get states of a node`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SetCPD`})}),Q(`td`,{children:Q(`code`,{children:`(name string, cpd *factors.TabularCPD)`})}),Q(`td`,{children:`Assign a CPD to a node`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`GetCPD`})}),Q(`td`,{children:Q(`code`,{children:`(name string) *factors.TabularCPD`})}),Q(`td`,{children:`Get the CPD of a node`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`CheckModel`})}),Q(`td`,{children:Q(`code`,{children:`() error`})}),Q(`td`,{children:`Validate DAG, CPDs, dimensions, probability sums`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ToMarkovFactors`})}),Q(`td`,{children:Q(`code`,{children:`() ([]*factors.DiscreteFactor, error)`})}),Q(`td`,{children:`Convert CPDs to factors for inference`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ToJunctionTree`})}),Q(`td`,{children:Q(`code`,{children:`() (*JunctionTree, error)`})}),Q(`td`,{children:`Build junction tree for BP inference`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DSeparation`})}),Q(`td`,{children:Q(`code`,{children:`(x, y, z []string) bool`})}),Q(`td`,{children:`Test d-separation: X _||_ Y | Z`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`GetIndependencies`})}),Q(`td`,{children:Q(`code`,{children:`() *independencies.Independencies`})}),Q(`td`,{children:`List all conditional independencies`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`bn := models.NewBayesianNetwork()
bn.AddNode("A")
bn.AddNode("B")
bn.AddNode("C")
bn.AddEdge("A", "B")
bn.AddEdge("B", "C")

bn.SetStates("A", []string{"a0", "a1"})
bn.SetStates("B", []string{"b0", "b1"})
bn.SetStates("C", []string{"c0", "c1"})

bn.SetCPD("A", factors.NewTabularCPD("A", 2, []float64{0.4, 0.6}, nil, nil))
bn.SetCPD("B", factors.NewTabularCPD("B", 2, []float64{0.2, 0.8, 0.9, 0.1}, []string{"A"}, []int{2}))
bn.SetCPD("C", factors.NewTabularCPD("C", 2, []float64{0.3, 0.7, 0.6, 0.4}, []string{"B"}, []int{2}))

err := bn.CheckModel()
fmt.Println("Valid:", err == nil)`})})]}),Q(`section`,{class:`section`,id:`api-factors`,children:[Q(`h2`,{children:`factors -- Factor Representations`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/factors`})]}),Q(`p`,{children:`Factor types and operations for probabilistic models. Factors are the building blocks of inference.`}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiscreteFactor`})}),Q(`td`,{children:`General discrete factor over a set of variables. Supports product, marginalize, reduce, normalize.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`TabularCPD`})}),Q(`td`,{children:`Conditional probability distribution table. CPD values are stored in column-major order.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`JointProbabilityDistribution`})}),Q(`td`,{children:`Full joint distribution over a set of variables.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LinearGaussianCPD`})}),Q(`td`,{children:`Linear Gaussian conditional: child = sum(beta_i * parent_i) + N(mean, variance).`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NoisyOR`})}),Q(`td`,{children:`Noisy-OR parameterization for compact CPD representation with many binary parents.`})]})]})]}),Q(`h3`,{children:`Constructors`}),Q(`pre`,{children:Q(`code`,{children:`// DiscreteFactor: general factor over named variables
f := factors.NewDiscreteFactor(
    []string{"A", "B"},     // variable names
    []int{2, 3},            // cardinalities
    []float64{...},         // values (product of cardinalities)
)

// TabularCPD: conditional probability distribution
cpd := factors.NewTabularCPD(
    "B",                    // variable name
    2,                      // number of states
    []float64{0.3, 0.7, 0.8, 0.2},  // values in column-major order
    []string{"A"},          // parent names (nil for root nodes)
    []int{2},               // parent cardinalities (nil for root nodes)
)

// JointProbabilityDistribution
jpd := factors.NewJPD(
    []string{"A", "B"},
    []int{2, 2},
    []float64{0.1, 0.2, 0.3, 0.4},
)

// LinearGaussianCPD
lgcpd := factors.NewLinearGaussianCPD(
    "Y",                    // variable name
    0.0,                    // mean of noise term
    1.0,                    // variance of noise term
    []string{"X"},          // parent names
    []float64{0.5},         // regression coefficients (betas)
)

// NoisyOR
nor := factors.NewNoisyOR(
    "Effect",
    []string{"Cause1", "Cause2"},
    []float64{0.1, 0.3},   // inhibition probabilities
    0.01,                   // leak probability
)`})}),Q(`h3`,{children:`DiscreteFactor Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Product`})}),Q(`td`,{children:Q(`code`,{children:`(other *DiscreteFactor) *DiscreteFactor`})}),Q(`td`,{children:`Factor product (join)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Marginalize`})}),Q(`td`,{children:Q(`code`,{children:`(vars []string) *DiscreteFactor`})}),Q(`td`,{children:`Sum out variables`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Reduce`})}),Q(`td`,{children:Q(`code`,{children:`(evidence map[string]int) *DiscreteFactor`})}),Q(`td`,{children:`Fix variables to observed values`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Normalize`})}),Q(`td`,{children:Q(`code`,{children:`() *DiscreteFactor`})}),Q(`td`,{children:`Normalize values to sum to 1`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Variables`})}),Q(`td`,{children:Q(`code`,{children:`() []string`})}),Q(`td`,{children:`Get variable names`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Cardinality`})}),Q(`td`,{children:Q(`code`,{children:`() []int`})}),Q(`td`,{children:`Get cardinalities`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Values`})}),Q(`td`,{children:Q(`code`,{children:`() *numgo.NDArray`})}),Q(`td`,{children:`Get values as NDArray`})]})]})]}),Q(`h3`,{children:`Example: Factor Operations`}),Q(`pre`,{children:Q(`code`,{children:`// Create two factors
f1 := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
f2 := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.5, 0.6, 0.7, 0.8})

// Factor product: f1 * f2 -> factor over {A, B, C}
product := f1.Product(f2)
fmt.Println("Product variables:", product.Variables())  // [A, B, C]

// Marginalize out B: sum_B (f1 * f2) -> factor over {A, C}
marginal := product.Marginalize([]string{"B"})
fmt.Println("Marginalized variables:", marginal.Variables())  // [A, C]

// Reduce: fix A=0 -> factor over {B, C}
reduced := product.Reduce(map[string]int{"A": 0})
fmt.Println("Reduced values:", reduced.Values().Data())

// Normalize
normalized := marginal.Normalize()
fmt.Println("Sum after normalize:", normalized.Values().Sum())  // 1.0`})})]}),Q(`section`,{class:`section`,id:`api-inference`,children:[Q(`h2`,{children:`inference -- Inference Algorithms`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/inference`})]}),Q(`p`,{children:`7 inference algorithms for computing posterior probabilities and MAP assignments.`}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`VariableElimination`})}),Q(`td`,{children:`Exact inference via factor elimination. Default method.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BeliefPropagation`})}),Q(`td`,{children:`Exact inference via message passing on junction trees.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MPLP`})}),Q(`td`,{children:`Max-Product Linear Programming for MAP inference.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ApproxInference`})}),Q(`td`,{children:`Sampling-based approximate inference.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`CausalInference`})}),Q(`td`,{children:`Do-calculus interventional queries.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DBNInference`})}),Q(`td`,{children:`Inference over dynamic Bayesian networks.`})]})]})]}),Q(`h3`,{children:`Constructors`}),Q(`pre`,{children:Q(`code`,{children:`// Variable Elimination (from Markov factors)
facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)

// Belief Propagation (from junction tree)
jt, _ := bn.ToJunctionTree()
cliques := jt.Cliques()
separators := jt.SeparatorSets()
cliqueFactors := make(map[int][]*factors.DiscreteFactor)
for i, c := range cliques {
    fs := jt.GetCliqueFactors(c)
    if len(fs) > 0 { cliqueFactors[i] = fs }
}
bp := inference.NewBeliefPropagation(cliques, separators, cliqueFactors)

// MPLP
mplp := inference.NewMPLP(facs)

// Approximate Inference
approx, _ := inference.NewApproxInference(bn, 10000)

// Causal Inference
ci, _ := inference.NewCausalInference(bn)

// DBN Inference
dbnInf, _ := inference.NewDBNInference(dbn)`})}),Q(`h3`,{children:`VariableElimination Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Query`})}),Q(`td`,{children:Q(`code`,{children:`(variables []string, evidence map[string]int) (*factors.DiscreteFactor, error)`})}),Q(`td`,{children:`Compute posterior P(variables | evidence)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MAP`})}),Q(`td`,{children:Q(`code`,{children:`(variables []string, evidence map[string]int) (map[string]int, error)`})}),Q(`td`,{children:`Find most likely assignment`})]})]})]}),Q(`h3`,{children:`CausalInference Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Query`})}),Q(`td`,{children:Q(`code`,{children:`(variables []string, doEvidence map[string]int, obsEvidence map[string]int) (*factors.DiscreteFactor, error)`})}),Q(`td`,{children:`Compute P(variables | do(interventions), observations)`})]})})]}),Q(`h3`,{children:`Example: VE Query and MAP`}),Q(`pre`,{children:Q(`code`,{children:`facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)

// Posterior query
result, _ := ve.Query([]string{"C"}, map[string]int{"A": 1})
fmt.Println("P(C | A=1):", result.Values().Data())

// MAP query
assignment, _ := ve.MAP([]string{"B"}, map[string]int{"A": 0})
fmt.Println("MAP(B | A=0):", assignment)`})}),Q(`h3`,{children:`Example: Causal Inference`}),Q(`pre`,{children:Q(`code`,{children:`ci, _ := inference.NewCausalInference(bn)

// P(Y | do(X=1))
result, _ := ci.Query([]string{"Y"}, map[string]int{"X": 1}, nil)
fmt.Println("P(Y | do(X=1)):", result.Values().Data())

// P(Y | do(X=1), Z=0)
result, _ = ci.Query([]string{"Y"}, map[string]int{"X": 1}, map[string]int{"Z": 0})`})})]}),Q(`section`,{class:`section`,id:`api-sampling`,children:[Q(`h2`,{children:`sampling -- Sampling Methods`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/sampling`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BayesianModelSampling`})}),Q(`td`,{children:`Forward, rejection, and likelihood-weighted sampling from a BN.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`GibbsSampling`})}),Q(`td`,{children:`MCMC Gibbs sampler with burn-in and thinning.`})]})]})]}),Q(`h3`,{children:`Constructors`}),Q(`pre`,{children:Q(`code`,{children:`bms, err := sampling.NewBayesianModelSampling(bn, 42)  // seed=42
gs, err := sampling.NewGibbsSampling(bn, 42)`})}),Q(`h3`,{children:`BayesianModelSampling Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ForwardSample`})}),Q(`td`,{children:Q(`code`,{children:`(n int) (*tabgo.DataFrame, error)`})}),Q(`td`,{children:`Generate n forward samples`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`RejectionSample`})}),Q(`td`,{children:Q(`code`,{children:`(n int, evidence map[string]int) (*tabgo.DataFrame, error)`})}),Q(`td`,{children:`Generate n samples consistent with evidence`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LikelihoodWeightedSample`})}),Q(`td`,{children:Q(`code`,{children:`(n int, evidence map[string]int) (*tabgo.DataFrame, []float64, error)`})}),Q(`td`,{children:`Generate n weighted samples`})]})]})]}),Q(`h3`,{children:`GibbsSampling Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Sample`})}),Q(`td`,{children:Q(`code`,{children:`(nSamples, burnIn, thinning int, evidence map[string]int) (*tabgo.DataFrame, error)`})}),Q(`td`,{children:`Run Gibbs sampler`})]})})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`bms, _ := sampling.NewBayesianModelSampling(bn, 42)

// Forward sampling
samples, _ := bms.ForwardSample(1000)
fmt.Printf("Forward: %d samples\\n", samples.NRows())

// Rejection sampling with evidence
rejSamples, _ := bms.RejectionSample(500, map[string]int{"A": 1})

// Gibbs sampling
gs, _ := sampling.NewGibbsSampling(bn, 42)
gibbsSamples, _ := gs.Sample(1000, 100, 2, map[string]int{"A": 1})`})})]}),Q(`section`,{class:`section`,id:`api-learning`,children:[Q(`h2`,{children:`learning -- Learning Algorithms`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/learning`})]}),Q(`p`,{children:`15+ algorithms for parameter estimation and structure learning.`}),Q(`h3`,{children:`Parameter Learning`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MLE`})}),Q(`td`,{children:Q(`code`,{children:`NewMLE(bn, data)`})}),Q(`td`,{children:`Maximum Likelihood Estimation from complete data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BayesianEstimator`})}),Q(`td`,{children:Q(`code`,{children:`NewBayesianEstimator(bn, data, prior, ess)`})}),Q(`td`,{children:`Bayesian estimation with Dirichlet prior (BDeu)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`EM`})}),Q(`td`,{children:Q(`code`,{children:`NewEM(bn, data, latent, maxIter, tol)`})}),Q(`td`,{children:`Expectation-Maximization for incomplete data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LinearGaussianMLE`})}),Q(`td`,{children:Q(`code`,{children:`NewLinearGaussianMLE(lgbn, data)`})}),Q(`td`,{children:`MLE for linear Gaussian BNs`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MarginalEstimator`})}),Q(`td`,{children:Q(`code`,{children:`NewMarginalEstimator(bn, data)`})}),Q(`td`,{children:`Marginal likelihood estimation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MirrorDescent`})}),Q(`td`,{children:Q(`code`,{children:`NewMirrorDescent(bn, data, lr, maxIter)`})}),Q(`td`,{children:`Mirror descent optimization`})]})]})]}),Q(`h3`,{children:`Structure Learning`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Category`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`HillClimbSearch`})}),Q(`td`,{children:Q(`code`,{children:`NewHillClimbSearch(data, scoreFn)`})}),Q(`td`,{children:`Score-based`}),Q(`td`,{children:`Greedy hill-climbing with add/remove/reverse edge operations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`PC`})}),Q(`td`,{children:Q(`code`,{children:`NewPC(data, ciTest, significance)`})}),Q(`td`,{children:`Constraint-based`}),Q(`td`,{children:`PC algorithm using CI tests`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`GES`})}),Q(`td`,{children:Q(`code`,{children:`NewGES(data, scoreFn)`})}),Q(`td`,{children:`Score-based`}),Q(`td`,{children:`Greedy Equivalence Search over DAG equivalence classes`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ExhaustiveSearch`})}),Q(`td`,{children:Q(`code`,{children:`NewExhaustiveSearch(data, scoreFn)`})}),Q(`td`,{children:`Score-based`}),Q(`td`,{children:`Exhaustive enumeration (small networks only)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`TreeSearch`})}),Q(`td`,{children:Q(`code`,{children:`NewTreeSearch(data)`})}),Q(`td`,{children:`Score-based`}),Q(`td`,{children:`Chow-Liu tree-structured network learning`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MMHC`})}),Q(`td`,{children:Q(`code`,{children:`NewMMHC(data, ciTest, scoreFn, significance)`})}),Q(`td`,{children:`Hybrid`}),Q(`td`,{children:`Max-Min Hill-Climbing`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ExpertInLoop`})}),Q(`td`,{children:Q(`code`,{children:`NewExpertInLoop(data, scoreFn, client)`})}),Q(`td`,{children:`Interactive`}),Q(`td`,{children:`Expert/LLM-guided structure learning`})]})]})]}),Q(`h3`,{children:`Causal Estimation`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`IVEstimator`})}),Q(`td`,{children:Q(`code`,{children:`NewIVEstimator(data, instrument, treatment, outcome)`})}),Q(`td`,{children:`Instrumental variable estimation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SEMEstimator`})}),Q(`td`,{children:Q(`code`,{children:`NewSEMEstimator(sem, data)`})}),Q(`td`,{children:`SEM parameter estimation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LLMClient`})}),Q(`td`,{children:Q(`code`,{children:`NewLLMClient(endpoint, apiKey)`})}),Q(`td`,{children:`LLM client for expert-in-the-loop`})]})]})]}),Q(`h3`,{children:`Common Methods`}),Q(`p`,{children:[`All parameter learners have `,Q(`code`,{children:`Estimate() error`}),` which populates CPDs on the model.`]}),Q(`p`,{children:[`All structure learners have `,Q(`code`,{children:`Estimate() (*models.BayesianNetwork, error)`}),` which returns a learned structure.`]}),Q(`p`,{children:[`PC also provides `,Q(`code`,{children:`EstimateBN() (*models.BayesianNetwork, error)`}),` and `,Q(`code`,{children:`EstimatePDAG() (*graphgo.PDAG, error)`}),`.`]}),Q(`h3`,{children:`Example: Complete Learning Pipeline`}),Q(`pre`,{children:Q(`code`,{children:`// Load data
data, _ := tabgo.ReadCSV("observations.csv")

// Learn structure
scorer := structure_score.NewBIC()
hc := learning.NewHillClimbSearch(data, scorer.LocalScore)
bn, _ := hc.Estimate()

// Fit parameters
mle := learning.NewMLE(bn, data)
mle.Estimate()

// Use the learned model
facs, _ := bn.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)
result, _ := ve.Query([]string{"Y"}, map[string]int{"X": 1})`})})]}),Q(`section`,{class:`section`,id:`api-ci-tests`,children:[Q(`h2`,{children:`ci_tests -- Conditional Independence Tests`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/ci_tests`})]}),Q(`p`,{children:`16 conditional independence tests used by constraint-based structure learning algorithms.`}),Q(`h3`,{children:`Tests by Category`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Category`}),Q(`th`,{children:`Function`}),Q(`th`,{children:`Data Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`Discrete`}),Q(`td`,{children:Q(`code`,{children:`ChiSquare`})}),Q(`td`,{children:`Categorical`}),Q(`td`,{children:`Pearson's chi-squared test`})]}),Q(`tr`,{children:[Q(`td`,{children:`Discrete`}),Q(`td`,{children:Q(`code`,{children:`GSq`})}),Q(`td`,{children:`Categorical`}),Q(`td`,{children:`G-squared (log-likelihood ratio) test`})]}),Q(`tr`,{children:[Q(`td`,{children:`Discrete`}),Q(`td`,{children:Q(`code`,{children:`ModifiedChiSquare`})}),Q(`td`,{children:`Categorical`}),Q(`td`,{children:`Modified chi-squared with correction`})]}),Q(`tr`,{children:[Q(`td`,{children:`Discrete`}),Q(`td`,{children:Q(`code`,{children:`PowerDivergence`})}),Q(`td`,{children:`Categorical`}),Q(`td`,{children:`Power divergence family (generalizes chi-sq and G-sq)`})]}),Q(`tr`,{children:[Q(`td`,{children:`Continuous`}),Q(`td`,{children:Q(`code`,{children:`FisherZ`})}),Q(`td`,{children:`Continuous`}),Q(`td`,{children:`Fisher's Z-transform of partial correlation`})]}),Q(`tr`,{children:[Q(`td`,{children:`Continuous`}),Q(`td`,{children:Q(`code`,{children:`Pearsonr`})}),Q(`td`,{children:`Continuous`}),Q(`td`,{children:`Pearson correlation test`})]}),Q(`tr`,{children:[Q(`td`,{children:`Continuous`}),Q(`td`,{children:Q(`code`,{children:`PartialCorrelation`})}),Q(`td`,{children:`Continuous`}),Q(`td`,{children:`Partial correlation test`})]}),Q(`tr`,{children:[Q(`td`,{children:`Multivariate`}),Q(`td`,{children:Q(`code`,{children:`GCM`})}),Q(`td`,{children:`Any`}),Q(`td`,{children:`Generalized Covariance Measure`})]}),Q(`tr`,{children:[Q(`td`,{children:`Multivariate`}),Q(`td`,{children:Q(`code`,{children:`HotellingLawley`})}),Q(`td`,{children:`Multivariate`}),Q(`td`,{children:`Hotelling-Lawley trace test`})]}),Q(`tr`,{children:[Q(`td`,{children:`Multivariate`}),Q(`td`,{children:Q(`code`,{children:`PillaiBartlett`})}),Q(`td`,{children:`Multivariate`}),Q(`td`,{children:`Pillai-Bartlett trace test`})]}),Q(`tr`,{children:[Q(`td`,{children:`Tree-based`}),Q(`td`,{children:Q(`code`,{children:`TreeBasedCI`})}),Q(`td`,{children:`Any`}),Q(`td`,{children:`Decision-tree-based CI test`})]})]})]}),Q(`h3`,{children:`Common Signature`}),Q(`pre`,{children:Q(`code`,{children:`func ChiSquare(x, y string, z []string, data *tabgo.DataFrame, significance float64) (statistic float64, pValue float64, independent bool)`})}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`stat, pValue, indep := ci_tests.ChiSquare("X", "Y", []string{"Z"}, data, 0.05)
fmt.Printf("Chi-square=%.3f, p=%.4f, independent=%v\\n", stat, pValue, indep)

stat2, pValue2, indep2 := ci_tests.FisherZ("X", "Y", []string{"Z"}, data, 0.05)
fmt.Printf("FisherZ=%.3f, p=%.4f, independent=%v\\n", stat2, pValue2, indep2)`})})]}),Q(`section`,{class:`section`,id:`api-structure-score`,children:[Q(`h2`,{children:`structure_score -- Scoring Functions`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/structure_score`})]}),Q(`p`,{children:`13 scoring function variants for score-based structure learning.`}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BIC`})}),Q(`td`,{children:Q(`code`,{children:`NewBIC()`})}),Q(`td`,{children:`Bayesian Information Criterion. Penalizes by log(N)/2.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`AIC`})}),Q(`td`,{children:Q(`code`,{children:`NewAIC()`})}),Q(`td`,{children:`Akaike Information Criterion. Less penalty than BIC.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BDeu`})}),Q(`td`,{children:Q(`code`,{children:`NewBDeu(equivalentSampleSize)`})}),Q(`td`,{children:`Bayesian Dirichlet equivalent uniform.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`BDs`})}),Q(`td`,{children:Q(`code`,{children:`NewBDs()`})}),Q(`td`,{children:`Bayesian Dirichlet sparse.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`K2`})}),Q(`td`,{children:Q(`code`,{children:`NewK2()`})}),Q(`td`,{children:`K2 score with uniform prior.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LogLikelihood`})}),Q(`td`,{children:Q(`code`,{children:`NewLogLikelihood()`})}),Q(`td`,{children:`Raw log-likelihood (no penalty).`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Gaussian`})}),Q(`td`,{children:Q(`code`,{children:`NewGaussian()`})}),Q(`td`,{children:`Score for continuous Gaussian data.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ConditionalGaussian`})}),Q(`td`,{children:Q(`code`,{children:`NewConditionalGaussian()`})}),Q(`td`,{children:`Score for mixed discrete/continuous.`})]})]})]}),Q(`h3`,{children:`Common Method`}),Q(`pre`,{children:Q(`code`,{children:`// All scorers provide LocalScore for use with structure learning
scorer := structure_score.NewBIC()
score := scorer.LocalScore(variable string, parents []string, data *tabgo.DataFrame) float64`})}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`bic := structure_score.NewBIC()
bdeu := structure_score.NewBDeu(5.0)  // equivalent sample size = 5

// Use with hill-climbing
hc := learning.NewHillClimbSearch(data, bic.LocalScore)
bn, _ := hc.Estimate()`})})]}),Q(`section`,{class:`section`,id:`api-identification`,children:[Q(`h2`,{children:`identification -- Causal Effect Identification`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/identification`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Adjustment`})}),Q(`td`,{children:Q(`code`,{children:`NewAdjustment(bn)`})}),Q(`td`,{children:`Back-door criterion. Finds valid adjustment sets to identify causal effects.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Frontdoor`})}),Q(`td`,{children:Q(`code`,{children:`NewFrontdoor(bn)`})}),Q(`td`,{children:`Front-door criterion. Identifies effects via mediating variables.`})]})]})]}),Q(`h3`,{children:`Methods`}),Q(`pre`,{children:Q(`code`,{children:`adj := identification.NewAdjustment(bn)
adjustmentSet, err := adj.GetAdjustmentSet("Treatment", "Outcome")
// Returns the set of variables to condition on for back-door adjustment

fd := identification.NewFrontdoor(bn)
frontdoorSet, err := fd.GetFrontdoorSet("Treatment", "Outcome")
// Returns the set of mediating variables for front-door adjustment`})})]}),Q(`section`,{class:`section`,id:`api-prediction`,children:[Q(`h2`,{children:`prediction -- Causal Prediction`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/prediction`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DoubleMLRegressor`})}),Q(`td`,{children:Q(`code`,{children:`NewDoubleMLRegressor(data, treatment, outcome, confounders)`})}),Q(`td`,{children:`Double/Debiased Machine Learning for ATE estimation. Handles high-dimensional confounders.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NaiveAdjustmentRegressor`})}),Q(`td`,{children:Q(`code`,{children:`NewNaiveAdjustmentRegressor(data, treatment, outcome, adjustmentSet)`})}),Q(`td`,{children:`Naive back-door adjustment regression.`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NaiveIVRegressor`})}),Q(`td`,{children:Q(`code`,{children:`NewNaiveIVRegressor(data, instrument, treatment, outcome)`})}),Q(`td`,{children:`Instrumental variable regression for causal effects with unmeasured confounding.`})]})]})]}),Q(`h3`,{children:`Common Methods`}),Q(`pre`,{children:Q(`code`,{children:`regressor := prediction.NewDoubleMLRegressor(data, "Treatment", "Outcome", []string{"X1", "X2"})
ate, err := regressor.EstimateATE()
fmt.Println("Average Treatment Effect:", ate)`})})]}),Q(`section`,{class:`section`,id:`api-metrics`,children:[Q(`h2`,{children:`metrics -- Model Evaluation`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/metrics`})]}),Q(`h3`,{children:`Functions`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Function`}),Q(`th`,{children:`Signature`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SHD`})}),Q(`td`,{children:Q(`code`,{children:`(true, est *graphgo.DiGraph) int`})}),Q(`td`,{children:`Structural Hamming Distance (edge additions + deletions + reversals)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`AdjacencyConfusionMatrix`})}),Q(`td`,{children:Q(`code`,{children:`(true, est *graphgo.DiGraph) (tp, fp, tn, fn int)`})}),Q(`td`,{children:`Edge presence TP/FP/TN/FN (ignoring direction)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`OrientationConfusionMatrix`})}),Q(`td`,{children:Q(`code`,{children:`(true, est *graphgo.DiGraph) (tp, fp, tn, fn int)`})}),Q(`td`,{children:`Edge orientation TP/FP/TN/FN`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`CorrelationScore`})}),Q(`td`,{children:Q(`code`,{children:`(true, est *graphgo.DiGraph) float64`})}),Q(`td`,{children:`Correlation-based structure similarity`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`FisherC`})}),Q(`td`,{children:Q(`code`,{children:`(bn *models.BayesianNetwork, data *tabgo.DataFrame) float64`})}),Q(`td`,{children:`Fisher's C statistic for model fit`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`LogLikelihoodScore`})}),Q(`td`,{children:Q(`code`,{children:`(bn *models.BayesianNetwork, data *tabgo.DataFrame) float64`})}),Q(`td`,{children:`Log-likelihood of data given model`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NormalizedSHD`})}),Q(`td`,{children:Q(`code`,{children:`(true, est *graphgo.DiGraph) float64`})}),Q(`td`,{children:`SHD normalized by max possible SHD`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`shd := metrics.SHD(trueGraph, learnedGraph)
fmt.Println("SHD:", shd)

tp, fp, tn, fn := metrics.AdjacencyConfusionMatrix(trueGraph, learnedGraph)
precision := float64(tp) / float64(tp + fp)
recall := float64(tp) / float64(tp + fn)
fmt.Printf("Precision=%.3f, Recall=%.3f\\n", precision, recall)`})})]}),Q(`section`,{class:`section`,id:`api-readwrite`,children:[Q(`h2`,{children:`readwrite -- File Format I/O`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/readwrite`})]}),Q(`p`,{children:[`10 file format readers and writers. All read functions take `,Q(`code`,{children:`io.Reader`}),`, all write functions take `,Q(`code`,{children:`io.Writer`}),`.`]}),Q(`h3`,{children:`Functions`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Read`}),Q(`th`,{children:`Write`}),Q(`th`,{children:`Format`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadBIF(r io.Reader) (*models.BayesianNetwork, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteBIF(w io.Writer, bn *models.BayesianNetwork) error`})}),Q(`td`,{children:`BIF`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadXMLBIF(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteXMLBIF(w, bn) error`})}),Q(`td`,{children:`XMLBIF`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadNET(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteNET(w, bn) error`})}),Q(`td`,{children:`NET`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadUAI(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteUAI(w, bn) error`})}),Q(`td`,{children:`UAI`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadXDSL(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteXDSL(w, bn) error`})}),Q(`td`,{children:`XDSL`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadPomdpX(r) (*BN, error)`})}),Q(`td`,{children:`--`}),Q(`td`,{children:`PomdpX (read only)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadXBN(r) (*BN, error)`})}),Q(`td`,{children:`--`}),Q(`td`,{children:`XBN (read only)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadCSVModel(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteCSVModel(w, bn) error`})}),Q(`td`,{children:`CSV`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadJSONModel(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteJSONModel(w, bn) error`})}),Q(`td`,{children:`JSON`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ReadXMLModel(r) (*BN, error)`})}),Q(`td`,{children:Q(`code`,{children:`WriteXMLModel(w, bn) error`})}),Q(`td`,{children:`XML`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`// Read BIF
f, _ := os.Open("model.bif")
bn, _ := readwrite.ReadBIF(f)
f.Close()

// Write JSON
out, _ := os.Create("model.json")
readwrite.WriteJSONModel(out, bn)
out.Close()`})})]}),Q(`section`,{class:`section`,id:`api-base`,children:[Q(`h2`,{children:`base -- Foundational Graph Types`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/base`})]}),Q(`p`,{children:`7 graph types that serve as the structural foundation for all model types.`}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DAG`})}),Q(`td`,{children:Q(`code`,{children:`NewDAG()`})}),Q(`td`,{children:`Directed acyclic graph with cycle detection, topological sort, d-separation`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`PDAG`})}),Q(`td`,{children:Q(`code`,{children:`NewPDAG()`})}),Q(`td`,{children:`Partially directed acyclic graph (for equivalence classes)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`UndirectedGraph`})}),Q(`td`,{children:Q(`code`,{children:`NewUndirectedGraph()`})}),Q(`td`,{children:`Undirected graph for Markov networks`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ADMG`})}),Q(`td`,{children:Q(`code`,{children:`NewADMG()`})}),Q(`td`,{children:`Acyclic directed mixed graph (directed + bidirected edges)`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`MAG`})}),Q(`td`,{children:Q(`code`,{children:`NewMAG()`})}),Q(`td`,{children:`Maximal ancestral graph`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`SimpleCausalModel`})}),Q(`td`,{children:Q(`code`,{children:`NewSimpleCausalModel()`})}),Q(`td`,{children:`Basic causal model with intervention semantics`})]})]})]}),Q(`h3`,{children:`DAG Methods`}),Q(`pre`,{children:Q(`code`,{children:`dag := base.NewDAG()
dag.AddNode("X")
dag.AddNode("Y")
dag.AddNode("Z")
dag.AddEdge("X", "Y")
dag.AddEdge("Y", "Z")

sorted := dag.TopologicalSort()
hasCycle := dag.HasCycle()
separated := dag.DSeparation([]string{"X"}, []string{"Z"}, []string{"Y"})
ancestors := dag.Ancestors("Z")
descendants := dag.Descendants("X")`})})]}),Q(`section`,{class:`section`,id:`api-independencies`,children:[Q(`h2`,{children:`independencies -- Independence Assertions`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/independencies`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`IndependenceAssertion`})}),Q(`td`,{children:`A single assertion: X _||_ Y | Z. Fields: Event1, Event2, Event3 (conditioning).`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Independencies`})}),Q(`td`,{children:`A collection of IndependenceAssertions.`})]})]})]}),Q(`pre`,{children:Q(`code`,{children:`// Get all independencies from a BN
indeps := bn.GetIndependencies()
for _, assertion := range indeps.Assertions() {
    fmt.Printf("%v _||_ %v | %v\\n", assertion.Event1, assertion.Event2, assertion.Event3)
}`})})]}),Q(`section`,{class:`section`,id:`api-config`,children:[Q(`h2`,{children:`config -- Configuration`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/config`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Config`})}),Q(`td`,{children:`Global configuration for tolerances, defaults, and behavior tuning.`})]})})]}),Q(`pre`,{children:Q(`code`,{children:`cfg := config.Global()
// Configuration controls numerical tolerances, default methods,
// convergence thresholds, and logging verbosity.`})})]}),Q(`section`,{class:`section`,id:`api-utils`,children:[Q(`h2`,{children:`utils -- Shared Utilities`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/pgm/utils`})]}),Q(`p`,{children:`Shared parsing, optimization, and compatibility utilities. Includes helper functions for set operations, string parsing, numerical utilities, and other common operations used across packages.`})]}),Q(`section`,{class:`section`,id:`api-numgo`,children:[Q(`h2`,{children:`numgo -- N-Dimensional Arrays and BLAS (numpy equivalent)`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/numgo`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`NDArray`})}),Q(`td`,{children:Q(`code`,{children:`NewNDArray(shape []int)`})}),Q(`td`,{children:`N-dimensional array with shape, stride, element-wise ops`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Matrix`})}),Q(`td`,{children:[Q(`code`,{children:`NewMatrix(rows, cols int)`}),` / `,Q(`code`,{children:`NewMatrixFromData(rows, cols, data)`})]}),Q(`td`,{children:`2D matrix with linear algebra operations`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Vector`})}),Q(`td`,{children:Q(`code`,{children:`NewVector(data []float64)`})}),Q(`td`,{children:`1D vector with dot product, norm, element-wise ops`})]})]})]}),Q(`h3`,{children:`NDArray Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Shape() []int`})}),Q(`td`,{children:`Get array dimensions`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Fill(value float64)`})}),Q(`td`,{children:`Set all elements to value`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Sum() float64`})}),Q(`td`,{children:`Sum all elements`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Data() []float64`})}),Q(`td`,{children:`Get underlying flat data`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Get(indices ...int) float64`})}),Q(`td`,{children:`Get element at indices`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Set(value float64, indices ...int)`})}),Q(`td`,{children:`Set element at indices`})]})]})]}),Q(`h3`,{children:`Matrix Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Multiply(other *Matrix) *Matrix`})}),Q(`td`,{children:`Matrix multiplication`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Transpose() *Matrix`})}),Q(`td`,{children:`Matrix transpose`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Inverse() *Matrix`})}),Q(`td`,{children:`Matrix inverse`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Det() float64`})}),Q(`td`,{children:`Determinant`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Rows() int`})}),Q(`td`,{children:`Number of rows`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Cols() int`})}),Q(`td`,{children:`Number of columns`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Get(i, j int) float64`})}),Q(`td`,{children:`Get element`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Set(i, j int, v float64)`})}),Q(`td`,{children:`Set element`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Data() []float64`})}),Q(`td`,{children:`Get flat data`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`m := numgo.NewMatrixFromData(2, 2, []float64{1, 2, 3, 4})
det := m.Det()     // -2
inv := m.Inverse()
identity := m.Multiply(inv)  // should be [[1,0],[0,1]]

v := numgo.NewVector([]float64{1, 2, 3})
norm := v.Norm()   // sqrt(14)`})})]}),Q(`section`,{class:`section`,id:`api-scigo`,children:[Q(`h2`,{children:`scigo -- Statistical Computing (scipy equivalent)`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/scigo`})]}),Q(`h3`,{children:`Distributions`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Methods`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Normal`})}),Q(`td`,{children:Q(`code`,{children:`NewNormal(mean, std)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`}),`, `,Q(`code`,{children:`Sample`}),`, `,Q(`code`,{children:`Mean`}),`, `,Q(`code`,{children:`Variance`})]})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`ChiSquared`})}),Q(`td`,{children:Q(`code`,{children:`NewChiSquared(df)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`}),`, `,Q(`code`,{children:`Sample`})]})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Beta`})}),Q(`td`,{children:Q(`code`,{children:`NewBeta(alpha, beta)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`}),`, `,Q(`code`,{children:`Mean`})]})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Gamma`})}),Q(`td`,{children:Q(`code`,{children:`NewGamma(shape, rate)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`}),`, `,Q(`code`,{children:`Mean`})]})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`StudentT`})}),Q(`td`,{children:Q(`code`,{children:`NewStudentT(df)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`})]})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Uniform`})}),Q(`td`,{children:Q(`code`,{children:`NewUniform(low, high)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`}),`, `,Q(`code`,{children:`Sample`})]})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Exponential`})}),Q(`td`,{children:Q(`code`,{children:`NewExponential(rate)`})}),Q(`td`,{children:[Q(`code`,{children:`PDF`}),`, `,Q(`code`,{children:`CDF`}),`, `,Q(`code`,{children:`PPF`})]})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`n := scigo.NewNormal(0, 1)
fmt.Println("CDF(1.96):", n.CDF(1.96))   // ~0.975
fmt.Println("PPF(0.975):", n.PPF(0.975))  // ~1.96

chi2 := scigo.NewChiSquared(5)
pValue := 1.0 - chi2.CDF(11.07)  // ~0.05`})})]}),Q(`section`,{class:`section`,id:`api-graphgo`,children:[Q(`h2`,{children:`graphgo -- Graph Algorithms (networkx equivalent)`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/graphgo`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Constructor`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DiGraph`})}),Q(`td`,{children:Q(`code`,{children:`NewDiGraph()`})}),Q(`td`,{children:`Directed graph`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Graph`})}),Q(`td`,{children:Q(`code`,{children:`NewGraph()`})}),Q(`td`,{children:`Undirected graph`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`PDAG`})}),Q(`td`,{children:Q(`code`,{children:`NewPDAG()`})}),Q(`td`,{children:`Partially directed acyclic graph`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`g := graphgo.NewDiGraph()
g.AddNode("A")
g.AddNode("B")
g.AddNode("C")
g.AddEdge("A", "B")
g.AddEdge("B", "C")

fmt.Println("Topo sort:", g.TopologicalSort())   // [A, B, C]
fmt.Println("Children of A:", g.Successors("A"))  // [B]
fmt.Println("Has cycle:", g.HasCycle())           // false`})})]}),Q(`section`,{class:`section`,id:`api-tabgo`,children:[Q(`h2`,{children:`tabgo -- Tabular Data (pandas equivalent)`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/tabgo`})]}),Q(`h3`,{children:`Types`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Type`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`DataFrame`})}),Q(`td`,{children:`Tabular data with named columns, filtering, groupby`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Series`})}),Q(`td`,{children:`Single column with value counts, unique, statistics`})]})]})]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`df, _ := tabgo.ReadCSV("data.csv")
fmt.Println("Rows:", df.NRows())
fmt.Println("Columns:", df.Columns())

col := df.Column("Score")
fmt.Println("Unique:", col.Unique())
fmt.Println("Counts:", col.ValueCounts())

tabgo.WriteCSV(df, "output.csv")`})})]}),Q(`section`,{class:`section`,id:`api-gpu`,children:[Q(`h2`,{children:`gpu -- GPU Compute Backend`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/gpu`})]}),Q(`p`,{children:`Optional GPU acceleration for compute-intensive operations. Provides GPU-backed matrix operations and factor computations for large-scale inference, learning, and deep learning training.`})]}),Q(`section`,{class:`section`,id:`api-tensorflow`,children:[Q(`h2`,{children:`tensorflow -- Deep Learning`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/lib/tensorflow`})]}),Q(`p`,{children:[`TensorFlow-compatible deep learning in pure Go. Provides neural network construction, training, and inference capabilities. Builds on `,Q(`code`,{children:`numgo`}),` for tensor operations and `,Q(`code`,{children:`gpu`}),` for optional hardware acceleration.`]})]}),Q(`section`,{class:`section`,id:`api-example-models`,children:[Q(`h2`,{children:`example_models -- Built-in Models`}),Q(`p`,{children:[`Import: `,Q(`code`,{children:`github.com/asymmetric-effort/datascience/example_models`})]}),Q(`h3`,{children:`Functions`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Function`}),Q(`th`,{children:`Description`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`List() []string`})}),Q(`td`,{children:`List all 25 available model names`})]}),Q(`tr`,{children:[Q(`td`,{children:Q(`code`,{children:`Get(name string) (*models.BayesianNetwork, error)`})}),Q(`td`,{children:`Load a model by name`})]})]})]}),Q(`h3`,{children:`Available Models`}),Q(`p`,{children:[Q(`strong`,{children:`With CPDs (13):`}),` student, asia, alarm, cancer, watersprinkler, survey, montyhall, dogproblem, frauddetection, medicaldiagnosis, earthquake, visitasia, cointoss`]}),Q(`p`,{children:[Q(`strong`,{children:`Structure-only (12):`}),` sachs, child, insurance, alarmfull, water, mildew, barley, hailfinder, hepar2, win95pts, pathfinder, pigs`]}),Q(`h3`,{children:`Example`}),Q(`pre`,{children:Q(`code`,{children:`names := example_models.List()
fmt.Println("Available:", len(names), "models")

bn, err := example_models.Get("asia")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%d nodes, %d edges\\n", len(bn.Nodes()), len(bn.Edges()))`})})]})]})}function Qn(){return y({title:`Tutorials — datascience`,description:`Step-by-step tutorials for building Bayesian networks, inference, structure learning, causal inference, file formats, sampling, advanced models, and internal libraries with datascience.`,canonical:`https://datascience.asymmetric-effort.com/#/tutorials`}),Q(`div`,{class:`page`,children:[Q(`h1`,{children:`Tutorials`}),Q(`nav`,{class:`page-toc`,children:[Q(`strong`,{children:`Tutorials:`}),` `,Q($,{to:`tutorial-1`,children:`1. First Bayesian Network`}),` | `,Q($,{to:`tutorial-2`,children:`2. Probabilistic Inference`}),` | `,Q($,{to:`tutorial-3`,children:`3. Learning from Data`}),` | `,Q($,{to:`tutorial-4`,children:`4. Causal Inference`}),` | `,Q($,{to:`tutorial-5`,children:`5. File Formats`}),` | `,Q($,{to:`tutorial-6`,children:`6. Sampling`}),` | `,Q($,{to:`tutorial-7`,children:`7. Advanced Models`}),` | `,Q($,{to:`tutorial-8`,children:`8. Internal Libraries`})]}),Q(`section`,{class:`section`,id:`tutorial-1`,children:[Q(`h2`,{children:`Tutorial 1: Building Your First Bayesian Network`}),Q(`p`,{children:`This tutorial walks through creating a Bayesian network from scratch. You will create the classic "wet grass" model, add conditional probability distributions, validate the model, run an inference query, and save it to a file. By the end you will understand the core workflow for using datascience.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:`Go 1.21+ installed`}),Q(`li`,{children:[`datascience installed: `,Q(`code`,{children:`go get github.com/asymmetric-effort/datascience`})]})]}),Q(`h3`,{children:`Introduction`}),Q(`p`,{children:`A Bayesian network (BN) is a directed acyclic graph (DAG) where each node represents a random variable and each directed edge represents a conditional dependency. The key idea: each variable's probability distribution depends only on its parents in the graph. This allows compact representation of joint probability distributions.`}),Q(`p`,{children:`We will model the classic "wet grass" scenario. On a given day, it might be cloudy. Cloudy weather makes rain more likely and the sprinkler less likely. Both rain and the sprinkler can cause the grass to be wet.`}),Q(`h3`,{children:`Step 1: Create the Network Structure`}),Q(`p`,{children:`Start by creating an empty BayesianNetwork, then add nodes (random variables) and directed edges (causal relationships).`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
    "github.com/asymmetric-effort/datascience/lib/pgm/readwrite"
    "os"
)

func main() {
    // Create an empty Bayesian network
    bn := models.NewBayesianNetwork()

    // Add four random variables as nodes
    bn.AddNode("Cloudy")
    bn.AddNode("Sprinkler")
    bn.AddNode("Rain")
    bn.AddNode("WetGrass")

    // Add directed edges representing causal influence:
    // Cloudy weather affects both the sprinkler and rain
    bn.AddEdge("Cloudy", "Sprinkler")
    bn.AddEdge("Cloudy", "Rain")
    // Both sprinkler and rain can cause wet grass
    bn.AddEdge("Sprinkler", "WetGrass")
    bn.AddEdge("Rain", "WetGrass")

    fmt.Println("Nodes:", bn.Nodes())
    // Output: Nodes: [Cloudy Sprinkler Rain WetGrass]

    fmt.Println("Edges:", bn.Edges())
    // Output: Edges: [[Cloudy Sprinkler] [Cloudy Rain] [Sprinkler WetGrass] [Rain WetGrass]]`})}),Q(`p`,{children:`At this point we have the graph structure but no probability information. The network knows the causal relationships but cannot answer probabilistic queries yet.`}),Q(`h3`,{children:`Step 2: Define States`}),Q(`p`,{children:`Each node needs a list of possible states (discrete values it can take). States are represented as string slices. The order matters: state 0 is the first string, state 1 is the second, and so on.`}),Q(`pre`,{children:Q(`code`,{children:`    // Define the possible values for each variable
    bn.SetStates("Cloudy", []string{"clear", "cloudy"})     // 0=clear, 1=cloudy
    bn.SetStates("Sprinkler", []string{"off", "on"})        // 0=off, 1=on
    bn.SetStates("Rain", []string{"no", "yes"})             // 0=no, 1=yes
    bn.SetStates("WetGrass", []string{"dry", "wet"})        // 0=dry, 1=wet`})}),Q(`h3`,{children:`Step 3: Add Conditional Probability Distributions (CPDs)`}),Q(`p`,{children:`Each node needs a CPD that specifies P(node | parents). Root nodes (no parents) have marginal distributions. Child nodes have conditional distributions.`}),Q(`p`,{children:[Q(`code`,{children:`NewTabularCPD(variable, cardinality, values, parents, parentCardinalities)`}),` creates a CPD. The `,Q(`code`,{children:`values`}),` array is in column-major order: iterate over parent combinations first, then over the variable's states.`]}),Q(`pre`,{children:Q(`code`,{children:`    // P(Cloudy) -- root node, no parents
    // 50% chance of clear, 50% chance of cloudy
    bn.SetCPD("Cloudy", factors.NewTabularCPD(
        "Cloudy", 2,                    // variable name, number of states
        []float64{0.5, 0.5},            // P(clear)=0.5, P(cloudy)=0.5
        nil, nil,                        // no parents
    ))

    // P(Sprinkler | Cloudy)
    // When clear: sprinkler on 50% of the time
    // When cloudy: sprinkler on only 10% of the time
    //
    // Values layout (column-major):
    //                  Cloudy=clear  Cloudy=cloudy
    // Sprinkler=off:     0.5           0.9
    // Sprinkler=on:      0.5           0.1
    bn.SetCPD("Sprinkler", factors.NewTabularCPD(
        "Sprinkler", 2,
        []float64{
            0.5, 0.9,   // P(off|clear)=0.5, P(off|cloudy)=0.9
            0.5, 0.1,   // P(on|clear)=0.5,  P(on|cloudy)=0.1
        },
        []string{"Cloudy"}, []int{2},
    ))

    // P(Rain | Cloudy)
    // When clear: 80% no rain. When cloudy: 80% rain.
    bn.SetCPD("Rain", factors.NewTabularCPD(
        "Rain", 2,
        []float64{
            0.8, 0.2,   // P(no|clear)=0.8, P(no|cloudy)=0.2
            0.2, 0.8,   // P(yes|clear)=0.2, P(yes|cloudy)=0.8
        },
        []string{"Cloudy"}, []int{2},
    ))

    // P(WetGrass | Sprinkler, Rain)
    // Two parents, each with 2 states = 4 columns
    // Column order: (S=off,R=no), (S=off,R=yes), (S=on,R=no), (S=on,R=yes)
    bn.SetCPD("WetGrass", factors.NewTabularCPD(
        "WetGrass", 2,
        []float64{
            // S=off,R=no  S=off,R=yes  S=on,R=no  S=on,R=yes
            1.0,          0.1,         0.1,       0.01,   // P(dry|...)
            0.0,          0.9,         0.9,       0.99,   // P(wet|...)
        },
        []string{"Sprinkler", "Rain"}, []int{2, 2},
    ))`})}),Q(`h3`,{children:`Step 4: Validate the Model`}),Q(`p`,{children:[Q(`code`,{children:`CheckModel()`}),` verifies that the graph is a valid DAG, every node has a CPD, CPD dimensions match the graph structure, and all probability distributions sum to 1.`]}),Q(`pre`,{children:Q(`code`,{children:`    // Validate the complete model
    if err := bn.CheckModel(); err != nil {
        log.Fatal("Model validation failed:", err)
    }
    fmt.Println("Model is valid!")
    // Output: Model is valid!`})}),Q(`h3`,{children:`Step 5: Run an Inference Query`}),Q(`p`,{children:`Convert the CPDs to Markov factors, create a Variable Elimination engine, and query posterior probabilities.`}),Q(`pre`,{children:Q(`code`,{children:`    // Convert CPDs to Markov factors for Variable Elimination
    facs, err := bn.ToMarkovFactors()
    if err != nil {
        log.Fatal(err)
    }
    ve := inference.NewVariableElimination(facs)

    // Query: P(Rain | WetGrass=wet)
    // "The grass is wet. Did it rain?"
    result, err := ve.Query(
        []string{"Rain"},                    // query variable
        map[string]int{"WetGrass": 1},       // evidence: WetGrass=1 means "wet"
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Rain | WetGrass=wet):", result.Values().Data())
    // Rain is more likely when grass is wet

    // MAP query: most likely assignment for Sprinkler given it rained
    assignment, err := ve.MAP(
        []string{"Sprinkler"},
        map[string]int{"Rain": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("MAP(Sprinkler | Rain=yes):", assignment)
    // Sprinkler is likely off when it's raining (both caused by cloudy weather)`})}),Q(`h3`,{children:`Step 6: Save the Model`}),Q(`pre`,{children:Q(`code`,{children:`    // Save the model as a BIF file
    f, err := os.Create("wetgrass.bif")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    if err := readwrite.WriteBIF(f, bn); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Model saved to wetgrass.bif")
}`})}),Q(`h3`,{children:`Expected Output`}),Q(`pre`,{children:Q(`code`,{children:`Nodes: [Cloudy Sprinkler Rain WetGrass]
Edges: [[Cloudy Sprinkler] [Cloudy Rain] [Sprinkler WetGrass] [Rain WetGrass]]
Model is valid!
P(Rain | WetGrass=wet): [0.2921... 0.7079...]
MAP(Sprinkler | Rain=yes): map[Sprinkler:0]
Model saved to wetgrass.bif`})}),Q(`h3`,{children:`What's Next`}),Q(`p`,{children:`Now that you can build and query a Bayesian network, explore:`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-2`,children:`Tutorial 2`}),`: Different inference methods (VE, BP, approximate)`]}),Q(`li`,{children:[Q($,{to:`tutorial-3`,children:`Tutorial 3`}),`: Learning structure and parameters from data`]}),Q(`li`,{children:[Q($,{to:`tutorial-4`,children:`Tutorial 4`}),`: Causal inference with do-calculus`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-2`,children:[Q(`h2`,{children:`Tutorial 2: Probabilistic Inference`}),Q(`p`,{children:`This tutorial covers the inference algorithms in datascience: Variable Elimination (exact), Belief Propagation (exact), and Approximate Inference (sampling-based). You will learn how to compute posterior marginals, MAP assignments, and when to use each method.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:[`Completed `,Q($,{to:`tutorial-1`,children:`Tutorial 1`})]}),Q(`li`,{children:`Understanding of conditional probability and Bayes' theorem`})]}),Q(`h3`,{children:`Introduction`}),Q(`p`,{children:`Inference is the process of computing posterior probabilities given a model and evidence. There are two main types:`}),Q(`ul`,{children:[Q(`li`,{children:[Q(`strong`,{children:`Marginal inference`}),`: P(query | evidence) -- what is the distribution of the query variables given what we observed?`]}),Q(`li`,{children:[Q(`strong`,{children:`MAP inference`}),`: argmax P(query | evidence) -- what is the most likely assignment of the query variables?`]})]}),Q(`h3`,{children:`Step 1: Load a Model`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

func main() {
    // Load the Asia (lung disease) model
    // 8 nodes: VisitAsia, Tuberculosis, Smoker, Lung, Bronc, TbOrCa, XRay, Dyspnea
    bn, err := example_models.Get("asia")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Loaded: %d nodes, %d edges\\n", len(bn.Nodes()), len(bn.Edges()))`})}),Q(`h3`,{children:`Step 2: Variable Elimination`}),Q(`p`,{children:`Variable Elimination (VE) is an exact inference algorithm. It works by systematically eliminating (summing out) variables from the joint distribution. It is the default method and works well for small to medium networks.`}),Q(`pre`,{children:Q(`code`,{children:`    // Convert to factors
    facs, err := bn.ToMarkovFactors()
    if err != nil {
        log.Fatal(err)
    }
    ve := inference.NewVariableElimination(facs)

    // Marginal query: P(Dyspnea)
    // No evidence -- just the prior distribution
    result, err := ve.Query([]string{"Dyspnea"}, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Dyspnea):", result.Values().Data())

    // Posterior query: P(Dyspnea | Smoker=yes)
    result, err = ve.Query(
        []string{"Dyspnea"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Dyspnea | Smoker=yes):", result.Values().Data())

    // Multiple query variables: P(Lung, Bronc | Smoker=yes)
    result, err = ve.Query(
        []string{"Lung", "Bronc"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Lung, Bronc | Smoker=yes):", result.Values().Data())

    // Multiple evidence: P(Dyspnea | Smoker=yes, VisitAsia=yes)
    result, err = ve.Query(
        []string{"Dyspnea"},
        map[string]int{"Smoker": 1, "VisitAsia": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("P(Dyspnea | Smoker=yes, VisitAsia=yes):", result.Values().Data())`})}),Q(`h3`,{children:`Step 3: MAP Queries`}),Q(`p`,{children:`MAP (Maximum A Posteriori) queries find the most likely assignment of query variables given evidence. This answers "what is the single best explanation?"`}),Q(`pre`,{children:Q(`code`,{children:`    // MAP: most likely values of Lung and Bronc given smoking
    assignment, err := ve.MAP(
        []string{"Lung", "Bronc"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("MAP(Lung, Bronc | Smoker=yes):", assignment)
    // Returns a map like: map[Bronc:1 Lung:0]
    // Meaning: most likely Bronc=yes, Lung=no`})}),Q(`h3`,{children:`Step 4: Belief Propagation`}),Q(`p`,{children:`Belief Propagation (BP) is another exact inference algorithm that works on junction trees. It is more efficient than VE when you need to answer many queries on the same model, because calibration is done once and subsequent queries are fast.`}),Q(`pre`,{children:Q(`code`,{children:`    // Build junction tree
    jt, err := bn.ToJunctionTree()
    if err != nil {
        log.Fatal(err)
    }

    // Set up Belief Propagation
    cliques := jt.Cliques()
    separators := jt.SeparatorSets()
    cliqueFactors := make(map[int][]*factors.DiscreteFactor)
    for i, c := range cliques {
        fs := jt.GetCliqueFactors(c)
        if len(fs) > 0 {
            cliqueFactors[i] = fs
        }
    }
    bp := inference.NewBeliefPropagation(cliques, separators, cliqueFactors)

    // Calibrate (run message passing) -- do this once
    if err := bp.Calibrate(); err != nil {
        log.Fatal(err)
    }

    // Query (fast after calibration)
    bpResult, err := bp.Query(
        []string{"Dyspnea"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("BP P(Dyspnea | Smoker=yes):", bpResult.Values().Data())
    // Should match VE result`})}),Q(`h3`,{children:`Step 5: Approximate Inference`}),Q(`p`,{children:`For large networks where exact inference is too slow, use sampling-based approximate inference. This generates many samples from the model and estimates posteriors from the empirical distribution.`}),Q(`pre`,{children:Q(`code`,{children:`    // Approximate inference using likelihood-weighted sampling
    approx, err := inference.NewApproxInference(bn, 10000) // 10,000 samples
    if err != nil {
        log.Fatal(err)
    }

    approxResult, err := approx.Query(
        []string{"Dyspnea"},
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Approx P(Dyspnea | Smoker=yes):", approxResult.Values().Data())
    // Close to but not exactly equal to exact result`})}),Q(`h3`,{children:`Comparing Methods`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Type`}),Q(`th`,{children:`Speed`}),Q(`th`,{children:`Accuracy`}),Q(`th`,{children:`When to Use`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`Variable Elimination`}),Q(`td`,{children:`Exact`}),Q(`td`,{children:`Fast for small networks`}),Q(`td`,{children:`Exact`}),Q(`td`,{children:`Default choice. Networks with up to ~50 nodes.`})]}),Q(`tr`,{children:[Q(`td`,{children:`Belief Propagation`}),Q(`td`,{children:`Exact`}),Q(`td`,{children:`Fast for many queries`}),Q(`td`,{children:`Exact`}),Q(`td`,{children:`When answering many queries on the same model.`})]}),Q(`tr`,{children:[Q(`td`,{children:`MPLP`}),Q(`td`,{children:`Exact (MAP)`}),Q(`td`,{children:`Fast`}),Q(`td`,{children:`Exact for MAP`}),Q(`td`,{children:`MAP inference on large networks.`})]}),Q(`tr`,{children:[Q(`td`,{children:`Approximate`}),Q(`td`,{children:`Approximate`}),Q(`td`,{children:`Scales to large networks`}),Q(`td`,{children:`Approximate`}),Q(`td`,{children:`Networks too large for exact inference.`})]})]})]}),Q(`h3`,{children:`Expected Output`}),Q(`pre`,{children:Q(`code`,{children:`Loaded: 8 nodes, 8 edges
P(Dyspnea): [0.560... 0.440...]
P(Dyspnea | Smoker=yes): [0.304... 0.696...]
MAP(Lung, Bronc | Smoker=yes): map[Bronc:1 Lung:0]
BP P(Dyspnea | Smoker=yes): [0.304... 0.696...]
Approx P(Dyspnea | Smoker=yes): [~0.30 ~0.70]`})}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-3`,children:`Tutorial 3`}),`: Learn models from data instead of specifying them manually`]}),Q(`li`,{children:[Q($,{to:`tutorial-4`,children:`Tutorial 4`}),`: Causal inference -- observation vs. intervention`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-3`,children:[Q(`h2`,{children:`Tutorial 3: Learning from Data`}),Q(`p`,{children:`This tutorial covers the complete learning pipeline: loading data, learning network structure using score-based (HillClimb) and constraint-based (PC) methods, fitting parameters with MLE and Bayesian estimation, and evaluating the result.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:[`Completed `,Q($,{to:`tutorial-1`,children:`Tutorial 1`}),` and `,Q($,{to:`tutorial-2`,children:`Tutorial 2`})]}),Q(`li`,{children:`CSV data file with discrete observations`})]}),Q(`h3`,{children:`Introduction`}),Q(`p`,{children:`In practice, you often have data but no known model structure. datascience provides two families of structure learning algorithms:`}),Q(`ul`,{children:[Q(`li`,{children:[Q(`strong`,{children:`Score-based`}),`: Search the space of DAGs to maximize a scoring function (BIC, BDeu, K2). Includes HillClimb, GES, ExhaustiveSearch, TreeSearch.`]}),Q(`li`,{children:[Q(`strong`,{children:`Constraint-based`}),`: Use conditional independence tests to discover the structure. Includes PC.`]}),Q(`li`,{children:[Q(`strong`,{children:`Hybrid`}),`: Combine both approaches. Includes MMHC.`]})]}),Q(`h3`,{children:`Step 1: Generate or Load Training Data`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/tabgo"
    "github.com/asymmetric-effort/datascience/lib/pgm/ci_tests"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
    "github.com/asymmetric-effort/datascience/lib/pgm/learning"
    "github.com/asymmetric-effort/datascience/lib/pgm/metrics"
    "github.com/asymmetric-effort/datascience/lib/pgm/sampling"
    "github.com/asymmetric-effort/datascience/lib/pgm/structure_score"
    "github.com/asymmetric-effort/datascience/lib/graphgo"
)

func main() {
    // Load the ground-truth Asia model
    trueBN, _ := example_models.Get("asia")

    // Generate 5000 training samples
    bms, _ := sampling.NewBayesianModelSampling(trueBN, 42)
    data, _ := bms.ForwardSample(5000)
    fmt.Printf("Training data: %d rows, %d columns\\n", data.NRows(), len(data.Columns()))

    // Alternatively, load from CSV:
    // data, _ := tabgo.ReadCSV("observations.csv")`})}),Q(`h3`,{children:`Step 2: Score-Based Learning (Hill-Climb with BIC)`}),Q(`p`,{children:`Hill-climbing starts with an empty graph and iteratively adds, removes, or reverses edges to maximize the BIC score. It is the most commonly used structure learning method.`}),Q(`pre`,{children:Q(`code`,{children:`    // Choose a scoring function
    scorer := structure_score.NewBIC()

    // Run hill-climbing
    hc := learning.NewHillClimbSearch(data, scorer.LocalScore)
    hcBN, err := hc.Estimate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hill-Climb: %d nodes, %d edges\\n",
        len(hcBN.Nodes()), len(hcBN.Edges()))`})}),Q(`h3`,{children:`Step 3: Constraint-Based Learning (PC Algorithm)`}),Q(`p`,{children:`The PC algorithm uses conditional independence (CI) tests to discover the graph skeleton (undirected edges), then orients edges based on v-structures and orientation rules.`}),Q(`pre`,{children:Q(`code`,{children:`    // Wrap the CI test function for the learning API
    // Signature: (x, y string, z []string, data *DataFrame, significance float64) -> (stat, pValue, independent)
    ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
        return ci_tests.ChiSquare(x, y, z, d, sig)
    }

    // Run PC with significance level 0.05
    pc := learning.NewPC(data, ciTest, 0.05)
    pcBN, err := pc.EstimateBN()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("PC: %d nodes, %d edges\\n",
        len(pcBN.Nodes()), len(pcBN.Edges()))`})}),Q(`h3`,{children:`Step 4: Other Structure Learning Methods`}),Q(`pre`,{children:Q(`code`,{children:`    // GES (Greedy Equivalence Search)
    ges := learning.NewGES(data, scorer.LocalScore)
    gesPDAG, err := ges.Estimate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("GES: learned PDAG")

    // Tree search (Chow-Liu) -- learns tree-structured BN
    ts := learning.NewTreeSearch(data)
    treeBN, err := ts.Estimate()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Tree: %d nodes, %d edges\\n",
        len(treeBN.Nodes()), len(treeBN.Edges()))

    // Exhaustive search (small networks only, up to ~5 nodes)
    // es := learning.NewExhaustiveSearch(smallData, scorer.LocalScore)
    // optBN, _ := es.Estimate()`})}),Q(`h3`,{children:`Step 5: Fit Parameters (MLE)`}),Q(`p`,{children:`After learning the structure, fit conditional probability distributions using Maximum Likelihood Estimation.`}),Q(`pre`,{children:Q(`code`,{children:`    // MLE: compute CPD parameters by counting occurrences in data
    mle := learning.NewMLE(hcBN, data)
    if err := mle.Estimate(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("MLE parameters fitted")`})}),Q(`h3`,{children:`Step 6: Fit Parameters (Bayesian Estimation)`}),Q(`p`,{children:`Bayesian estimation adds a Dirichlet prior, which helps with small datasets or rare state combinations.`}),Q(`pre`,{children:Q(`code`,{children:`    // Bayesian estimation with BDeu prior, equivalent sample size = 5.0
    be := learning.NewBayesianEstimator(hcBN, data, learning.BDeu, 5.0)
    if err := be.Estimate(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Bayesian parameters fitted")`})}),Q(`h3`,{children:`Step 7: EM for Incomplete Data`}),Q(`pre`,{children:Q(`code`,{children:`    // If data has missing values, use Expectation-Maximization
    // em := learning.NewEM(hcBN, incompleteData, nil, 100, 1e-6)
    // if err := em.Estimate(); err != nil {
    //     log.Fatal(err)
    // }`})}),Q(`h3`,{children:`Step 8: Evaluate the Learned Structure`}),Q(`pre`,{children:Q(`code`,{children:`    // Convert BNs to DiGraphs for comparison
    toDigraph := func(bn interface{ Nodes() []string; Edges() [][2]string }) *graphgo.DiGraph {
        g := graphgo.NewDiGraph()
        for _, n := range bn.Nodes() { g.AddNode(n) }
        for _, e := range bn.Edges() { g.AddEdge(e[0], e[1]) }
        return g
    }

    trueG := toDigraph(trueBN)
    hcG := toDigraph(hcBN)
    pcG := toDigraph(pcBN)

    // Structural Hamming Distance (lower is better)
    fmt.Println("Hill-Climb SHD:", metrics.SHD(trueG, hcG))
    fmt.Println("PC SHD:", metrics.SHD(trueG, pcG))

    // Adjacency confusion matrix
    aTP, aFP, aTN, aFN := metrics.AdjacencyConfusionMatrix(trueG, hcG)
    fmt.Printf("HC Adjacency: TP=%d FP=%d TN=%d FN=%d\\n", aTP, aFP, aTN, aFN)

    // Orientation confusion matrix
    oTP, oFP, oTN, oFN := metrics.OrientationConfusionMatrix(trueG, hcG)
    fmt.Printf("HC Orientation: TP=%d FP=%d TN=%d FN=%d\\n", oTP, oFP, oTN, oFN)`})}),Q(`h3`,{children:`Step 9: Run Inference on the Learned Model`}),Q(`pre`,{children:Q(`code`,{children:`    // The learned model is now fully parameterized -- use it for inference
    facs, _ := hcBN.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    result, _ := ve.Query(
        []string{"Dyspnea"},
        map[string]int{"Smoker": 1},
    )
    fmt.Println("Learned model P(Dyspnea | Smoker=yes):", result.Values().Data())
}`})}),Q(`h3`,{children:`Choosing a Scoring Function`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Score`}),Q(`th`,{children:`Best For`}),Q(`th`,{children:`Notes`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`BIC`}),Q(`td`,{children:`General use`}),Q(`td`,{children:`Good balance of fit and complexity. Consistent (recovers true structure with enough data).`})]}),Q(`tr`,{children:[Q(`td`,{children:`BDeu`}),Q(`td`,{children:`Small datasets`}),Q(`td`,{children:`Bayesian prior helps when data is scarce. Set equivalent sample size carefully.`})]}),Q(`tr`,{children:[Q(`td`,{children:`K2`}),Q(`td`,{children:`Fast evaluation`}),Q(`td`,{children:`Simple uniform prior. Fast to compute.`})]}),Q(`tr`,{children:[Q(`td`,{children:`AIC`}),Q(`td`,{children:`Less penalization`}),Q(`td`,{children:`Tends to learn denser networks than BIC.`})]})]})]}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-4`,children:`Tutorial 4`}),`: Causal inference on the learned model`]}),Q(`li`,{children:[Q($,{to:`tutorial-6`,children:`Tutorial 6`}),`: Generating data via sampling`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-4`,children:[Q(`h2`,{children:`Tutorial 4: Causal Inference`}),Q(`p`,{children:`This tutorial covers causal reasoning with datascience. You will learn the difference between observational and interventional queries, use do-calculus, compute average treatment effects, and use back-door adjustment.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:[`Completed `,Q($,{to:`tutorial-1`,children:`Tutorial 1`}),` and `,Q($,{to:`tutorial-2`,children:`Tutorial 2`})]}),Q(`li`,{children:`Basic understanding of causation vs. correlation`})]}),Q(`h3`,{children:`Introduction`}),Q(`p`,{children:`Bayesian networks encode causal relationships through directed edges. This enables two types of reasoning:`}),Q(`ul`,{children:[Q(`li`,{children:[Q(`strong`,{children:`Observational`}),`: P(Y | X=x) -- "If I observe X=x, what do I expect for Y?" This reflects correlation and may be confounded.`]}),Q(`li`,{children:[Q(`strong`,{children:`Interventional`}),`: P(Y | do(X=x)) -- "If I force X to be x (intervene), what happens to Y?" This reflects causation. The do() operator "cuts" incoming edges to X, removing confounding.`]})]}),Q(`h3`,{children:`Step 1: Build a Causal Model`}),Q(`p`,{children:`Consider a model of smoking, tar deposits, and cancer. Smoking causes tar deposits, and both smoking and tar cause cancer.`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/lib/pgm/factors"
    "github.com/asymmetric-effort/datascience/lib/pgm/identification"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
)

func main() {
    bn := models.NewBayesianNetwork()
    bn.AddNode("Smoking")
    bn.AddNode("Tar")
    bn.AddNode("Cancer")

    // Causal structure:
    // Smoking -> Tar -> Cancer
    // Smoking -> Cancer (direct effect)
    bn.AddEdge("Smoking", "Tar")
    bn.AddEdge("Smoking", "Cancer")
    bn.AddEdge("Tar", "Cancer")

    bn.SetStates("Smoking", []string{"no", "yes"})
    bn.SetStates("Tar", []string{"low", "high"})
    bn.SetStates("Cancer", []string{"no", "yes"})

    bn.SetCPD("Smoking", factors.NewTabularCPD(
        "Smoking", 2, []float64{0.7, 0.3}, nil, nil,
    ))
    bn.SetCPD("Tar", factors.NewTabularCPD(
        "Tar", 2,
        []float64{0.95, 0.1, 0.05, 0.9},
        []string{"Smoking"}, []int{2},
    ))
    bn.SetCPD("Cancer", factors.NewTabularCPD(
        "Cancer", 2,
        []float64{
            0.98, 0.9, 0.7, 0.05,    // P(no cancer | ...)
            0.02, 0.1, 0.3, 0.95,    // P(cancer | ...)
        },
        []string{"Smoking", "Tar"}, []int{2, 2},
    ))
    bn.CheckModel()`})}),Q(`h3`,{children:`Step 2: Observational vs. Interventional Query`}),Q(`pre`,{children:Q(`code`,{children:`    // OBSERVATIONAL: P(Cancer | Smoking=yes)
    // "Among people who smoke, what is the cancer rate?"
    facs, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    obsResult, _ := ve.Query(
        []string{"Cancer"},
        map[string]int{"Smoking": 1},
    )
    fmt.Println("Observational P(Cancer | Smoking=yes):", obsResult.Values().Data())

    // INTERVENTIONAL: P(Cancer | do(Smoking=yes))
    // "If we forced everyone to smoke, what would the cancer rate be?"
    // The do() operator removes incoming edges to Smoking,
    // isolating the causal effect from confounders.
    ci, err := inference.NewCausalInference(bn)
    if err != nil {
        log.Fatal(err)
    }
    doResult, _ := ci.Query(
        []string{"Cancer"},
        map[string]int{"Smoking": 1},   // do(Smoking=yes)
        nil,                              // no additional evidence
    )
    fmt.Println("Interventional P(Cancer | do(Smoking=yes)):", doResult.Values().Data())
    // In this model, the two may differ because Smoking has
    // both direct and indirect effects on Cancer.`})}),Q(`h3`,{children:`Step 3: Computing Average Treatment Effect (ATE)`}),Q(`p`,{children:`The ATE measures the causal effect of a treatment (Smoking) on an outcome (Cancer): ATE = P(Cancer=yes | do(Smoking=yes)) - P(Cancer=yes | do(Smoking=no))`}),Q(`pre`,{children:Q(`code`,{children:`    // P(Cancer | do(Smoking=yes))
    doSmokeYes, _ := ci.Query(
        []string{"Cancer"},
        map[string]int{"Smoking": 1},
        nil,
    )

    // P(Cancer | do(Smoking=no))
    doSmokeNo, _ := ci.Query(
        []string{"Cancer"},
        map[string]int{"Smoking": 0},
        nil,
    )

    // ATE = P(Cancer=yes | do(Smoke=yes)) - P(Cancer=yes | do(Smoke=no))
    pCancerYesTreated := doSmokeYes.Values().Data()[1]
    pCancerYesControl := doSmokeNo.Values().Data()[1]
    ate := pCancerYesTreated - pCancerYesControl
    fmt.Printf("ATE of Smoking on Cancer: %.4f\\n", ate)
    // Positive ATE means smoking causes increased cancer risk`})}),Q(`h3`,{children:`Step 4: Interventional Query with Evidence`}),Q(`pre`,{children:Q(`code`,{children:`    // P(Cancer | do(Tar=high), Smoking=no)
    // "If we forced tar deposits to be high but the person doesn't smoke,
    //  what is the cancer risk?"
    result, _ := ci.Query(
        []string{"Cancer"},
        map[string]int{"Tar": 1},         // intervene on Tar
        map[string]int{"Smoking": 0},      // observe Smoking=no
    )
    fmt.Println("P(Cancer | do(Tar=high), Smoking=no):", result.Values().Data())`})}),Q(`h3`,{children:`Step 5: Causal Identification`}),Q(`p`,{children:`Before computing a causal effect, check whether it is identifiable (computable from observational data) using the back-door or front-door criterion.`}),Q(`pre`,{children:Q(`code`,{children:`    // Back-door criterion: find an adjustment set
    adj := identification.NewAdjustment(bn)
    adjustmentSet, err := adj.GetAdjustmentSet("Smoking", "Cancer")
    if err == nil {
        fmt.Println("Back-door adjustment set:", adjustmentSet)
        // You can use this set to compute causal effects from
        // observational data by adjusting (conditioning) on these variables.
    } else {
        fmt.Println("No valid back-door adjustment set found")
    }

    // Front-door criterion
    fd := identification.NewFrontdoor(bn)
    frontdoorSet, err := fd.GetFrontdoorSet("Smoking", "Cancer")
    if err == nil {
        fmt.Println("Front-door set:", frontdoorSet)
        // Tar mediates the effect and can be used for front-door adjustment
    } else {
        fmt.Println("No valid front-door set found")
    }
}`})}),Q(`h3`,{children:`The Intuition`}),Q(`p`,{children:[Q(`strong`,{children:`Why do observation and intervention differ?`}),` When you observe X=x, you are conditioning on X. This selects a subpopulation where X happens to be x, which may carry information about confounders. When you intervene (do(X=x)), you force X to be x for everyone, breaking the connection between X and its causes. This isolates the downstream causal effect.`]}),Q(`p`,{children:`In a simple chain A -> B -> C, observing B gives information about A (through the chain), while do(B) does not -- the intervention breaks the A->B link.`}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-5`,children:`Tutorial 5`}),`: Working with different file formats`]}),Q(`li`,{children:[Q($,{to:`tutorial-7`,children:`Tutorial 7`}),`: Advanced model types (SEMs, Markov networks)`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-5`,children:[Q(`h2`,{children:`Tutorial 5: Working with File Formats`}),Q(`p`,{children:`datascience supports 10 file formats for model serialization. This tutorial covers reading, writing, converting between formats, and using the CLI for format conversion.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:[`Completed `,Q($,{to:`tutorial-1`,children:`Tutorial 1`})]}),Q(`li`,{children:`A BIF file (or create one in Tutorial 1)`})]}),Q(`h3`,{children:`Step 1: Reading Models from Files`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"
    "os"

    "github.com/asymmetric-effort/datascience/lib/pgm/readwrite"
)

func main() {
    // Read a BIF file
    f, err := os.Open("asia.bif")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    bn, err := readwrite.ReadBIF(f)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Loaded from BIF: %d nodes, %d edges\\n",
        len(bn.Nodes()), len(bn.Edges()))

    // Read other formats -- same pattern
    // f2, _ := os.Open("model.xmlbif")
    // bn2, _ := readwrite.ReadXMLBIF(f2)
    //
    // f3, _ := os.Open("model.net")
    // bn3, _ := readwrite.ReadNET(f3)
    //
    // f4, _ := os.Open("model.uai")
    // bn4, _ := readwrite.ReadUAI(f4)
    //
    // f5, _ := os.Open("model.xdsl")
    // bn5, _ := readwrite.ReadXDSL(f5)`})}),Q(`h3`,{children:`Step 2: Writing Models to Files`}),Q(`pre`,{children:Q(`code`,{children:`    // Write as BIF
    bifOut, _ := os.Create("output.bif")
    readwrite.WriteBIF(bifOut, bn)
    bifOut.Close()

    // Write as XMLBIF
    xmlbifOut, _ := os.Create("output.xmlbif")
    readwrite.WriteXMLBIF(xmlbifOut, bn)
    xmlbifOut.Close()

    // Write as NET (Hugin format)
    netOut, _ := os.Create("output.net")
    readwrite.WriteNET(netOut, bn)
    netOut.Close()

    // Write as UAI (competition format)
    uaiOut, _ := os.Create("output.uai")
    readwrite.WriteUAI(uaiOut, bn)
    uaiOut.Close()

    // Write as XDSL (GeNIe format)
    xdslOut, _ := os.Create("output.xdsl")
    readwrite.WriteXDSL(xdslOut, bn)
    xdslOut.Close()`})}),Q(`h3`,{children:`Step 3: JSON for Web Applications`}),Q(`p`,{children:`JSON is ideal for REST APIs, web dashboards, and JavaScript interop.`}),Q(`pre`,{children:Q(`code`,{children:`    // Write model as JSON
    jsonOut, _ := os.Create("model.json")
    readwrite.WriteJSONModel(jsonOut, bn)
    jsonOut.Close()
    fmt.Println("Wrote model.json")

    // Read model from JSON
    jsonIn, _ := os.Open("model.json")
    bnFromJSON, _ := readwrite.ReadJSONModel(jsonIn)
    jsonIn.Close()
    fmt.Printf("Read from JSON: %d nodes\\n", len(bnFromJSON.Nodes()))`})}),Q(`h3`,{children:`Step 4: CSV for Data Exchange`}),Q(`pre`,{children:Q(`code`,{children:`    // CSV model serialization
    csvOut, _ := os.Create("model.csv")
    readwrite.WriteCSVModel(csvOut, bn)
    csvOut.Close()

    // XML model serialization
    xmlOut, _ := os.Create("model.xml")
    readwrite.WriteXMLModel(xmlOut, bn)
    xmlOut.Close()`})}),Q(`h3`,{children:`Step 5: Format Conversion Pipeline`}),Q(`pre`,{children:Q(`code`,{children:`    // Read NET, convert to BIF and JSON
    netFile, _ := os.Open("model.net")
    bnFromNET, _ := readwrite.ReadNET(netFile)
    netFile.Close()

    bifFile, _ := os.Create("converted.bif")
    readwrite.WriteBIF(bifFile, bnFromNET)
    bifFile.Close()

    jsonFile, _ := os.Create("converted.json")
    readwrite.WriteJSONModel(jsonFile, bnFromNET)
    jsonFile.Close()
    fmt.Println("Converted NET -> BIF and JSON")
}`})}),Q(`h3`,{children:`Step 6: Working with Built-in Example Models`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/example_models"

// List all 25 available models
for _, name := range example_models.List() {
    fmt.Println(name)
}

// Load and export a model to BIF
bn, _ := example_models.Get("alarm")
f, _ := os.Create("alarm.bif")
readwrite.WriteBIF(f, bn)
f.Close()

// Load and export to JSON for a web API
bn, _ = example_models.Get("asia")
jsonOut, _ := os.Create("asia.json")
readwrite.WriteJSONModel(jsonOut, bn)
jsonOut.Close()`})}),Q(`h3`,{children:`Format Comparison`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Format`}),Q(`th`,{children:`Human Readable`}),Q(`th`,{children:`Tool Support`}),Q(`th`,{children:`Best For`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`BIF`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Broad`}),Q(`td`,{children:`General exchange, version control`})]}),Q(`tr`,{children:[Q(`td`,{children:`XMLBIF`}),Q(`td`,{children:`Somewhat`}),Q(`td`,{children:`Broad`}),Q(`td`,{children:`XML-based toolchains`})]}),Q(`tr`,{children:[Q(`td`,{children:`NET`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Hugin`}),Q(`td`,{children:`Hugin software users`})]}),Q(`tr`,{children:[Q(`td`,{children:`UAI`}),Q(`td`,{children:`Minimal`}),Q(`td`,{children:`Competitions`}),Q(`td`,{children:`UAI inference competitions`})]}),Q(`tr`,{children:[Q(`td`,{children:`XDSL`}),Q(`td`,{children:`Somewhat`}),Q(`td`,{children:`GeNIe/SMILE`}),Q(`td`,{children:`GeNIe software users`})]}),Q(`tr`,{children:[Q(`td`,{children:`JSON`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Universal`}),Q(`td`,{children:`Web APIs, JavaScript interop`})]}),Q(`tr`,{children:[Q(`td`,{children:`CSV`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Universal`}),Q(`td`,{children:`Spreadsheets, simple tools`})]}),Q(`tr`,{children:[Q(`td`,{children:`XML`}),Q(`td`,{children:`Somewhat`}),Q(`td`,{children:`Universal`}),Q(`td`,{children:`Enterprise systems`})]})]})]}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-6`,children:`Tutorial 6`}),`: Sampling data from models`]}),Q(`li`,{children:[Q($,{to:`tutorial-8`,children:`Tutorial 8`}),`: Using internal libraries for data processing`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-6`,children:[Q(`h2`,{children:`Tutorial 6: Sampling and Monte Carlo Methods`}),Q(`p`,{children:`This tutorial covers all sampling methods in datascience: forward sampling, rejection sampling, likelihood-weighted sampling, and Gibbs sampling. You will learn when to use each method and how to compare empirical distributions against exact marginals.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:[`Completed `,Q($,{to:`tutorial-1`,children:`Tutorial 1`})]}),Q(`li`,{children:`Basic understanding of Monte Carlo methods`})]}),Q(`h3`,{children:`Introduction`}),Q(`p`,{children:`Sampling generates data points (samples) from the joint distribution defined by a Bayesian network. Uses include:`}),Q(`ul`,{children:[Q(`li`,{children:`Generating synthetic training data`}),Q(`li`,{children:`Approximate inference when exact methods are too expensive`}),Q(`li`,{children:`Validating models by comparing sampled distributions to expected distributions`}),Q(`li`,{children:`Monte Carlo estimation of probabilities and expectations`})]}),Q(`h3`,{children:`Step 1: Forward Sampling`}),Q(`p`,{children:`Forward sampling generates samples by traversing the DAG in topological order. Each variable is sampled from its CPD given its parents' values. This is the simplest and most efficient method but cannot incorporate evidence.`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "log"

    "github.com/asymmetric-effort/datascience/example_models"
    "github.com/asymmetric-effort/datascience/lib/pgm/sampling"
    "github.com/asymmetric-effort/datascience/lib/tabgo"
)

func main() {
    bn, _ := example_models.Get("asia")

    // Create sampler with seed=42 for reproducibility
    bms, err := sampling.NewBayesianModelSampling(bn, 42)
    if err != nil {
        log.Fatal(err)
    }

    // Generate 10,000 forward samples
    samples, err := bms.ForwardSample(10000)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Generated %d samples with %d columns\\n",
        samples.NRows(), len(samples.Columns()))

    // Estimate P(Smoker=yes) from samples
    smokerCol := samples.Column("Smoker")
    counts := smokerCol.ValueCounts()
    fmt.Println("Smoker value counts:", counts)
    // Should be approximately 50/50 (depends on the model)

    // Save to CSV
    tabgo.WriteCSV(samples, "asia_samples.csv")
    fmt.Println("Saved to asia_samples.csv")`})}),Q(`h3`,{children:`Step 2: Rejection Sampling`}),Q(`p`,{children:`Rejection sampling incorporates evidence by generating forward samples and keeping only those consistent with the evidence. Simple but inefficient when evidence is rare.`}),Q(`pre`,{children:Q(`code`,{children:`    // Rejection sampling: generate samples consistent with Smoker=yes
    rejSamples, err := bms.RejectionSample(1000, map[string]int{"Smoker": 1})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Rejection samples: %d (all have Smoker=1)\\n", rejSamples.NRows())

    // Estimate P(Dyspnea=yes | Smoker=yes) from rejection samples
    dyspneaCol := rejSamples.Column("Dyspnea")
    dyspneaCounts := dyspneaCol.ValueCounts()
    fmt.Println("Dyspnea | Smoker=yes:", dyspneaCounts)`})}),Q(`h3`,{children:`Step 3: Likelihood-Weighted Sampling`}),Q(`p`,{children:`Likelihood-weighted sampling is more efficient than rejection sampling. Instead of discarding samples, it assigns weights to each sample based on the likelihood of the evidence. All samples are used, but weighted.`}),Q(`pre`,{children:Q(`code`,{children:`    // Likelihood-weighted sampling with evidence
    lwSamples, weights, err := bms.LikelihoodWeightedSample(
        5000,
        map[string]int{"Smoker": 1},
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("LW samples: %d with weights\\n", lwSamples.NRows())

    // The weights account for the evidence
    // Use weighted counts to estimate posteriors
    _ = weights`})}),Q(`h3`,{children:`Step 4: Gibbs Sampling (MCMC)`}),Q(`p`,{children:`Gibbs sampling is a Markov Chain Monte Carlo (MCMC) method. It starts with an initial assignment and iteratively resamples each variable from its full conditional distribution. It can handle evidence and is effective for large networks.`}),Q(`pre`,{children:Q(`code`,{children:`    // Create Gibbs sampler
    gs, err := sampling.NewGibbsSampling(bn, 42)
    if err != nil {
        log.Fatal(err)
    }

    // Parameters:
    //   nSamples=2000 -- number of samples to collect
    //   burnIn=500    -- discard first 500 samples (let chain converge)
    //   thinning=2    -- keep every 2nd sample (reduce autocorrelation)
    //   evidence      -- fix these variables
    gibbsSamples, err := gs.Sample(
        2000,                            // number of samples
        500,                             // burn-in
        2,                               // thinning interval
        map[string]int{"Smoker": 1},     // evidence
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Gibbs samples: %d\\n", gibbsSamples.NRows())

    // Estimate P(Dyspnea | Smoker=yes) from Gibbs samples
    gibbsDyspnea := gibbsSamples.Column("Dyspnea")
    fmt.Println("Gibbs Dyspnea | Smoker=yes:", gibbsDyspnea.ValueCounts())`})}),Q(`h3`,{children:`Step 5: Compare Empirical vs. Exact Marginals`}),Q(`pre`,{children:Q(`code`,{children:`    import "github.com/asymmetric-effort/datascience/lib/pgm/inference"

    // Exact posterior: P(Dyspnea | Smoker=yes)
    facs, _ := bn.ToMarkovFactors()
    ve := inference.NewVariableElimination(facs)
    exact, _ := ve.Query([]string{"Dyspnea"}, map[string]int{"Smoker": 1})
    fmt.Println("Exact P(Dyspnea | Smoker=yes):", exact.Values().Data())

    // The sampling-based estimates should converge to the exact values
    // as the number of samples increases.
    fmt.Println("\\nComparison:")
    fmt.Println("  Exact:     ", exact.Values().Data())
    fmt.Println("  Rejection: (from counts above)")
    fmt.Println("  Gibbs:     (from counts above)")
}`})}),Q(`h3`,{children:`When to Use Each Method`}),Q(`table`,{children:[Q(`thead`,{children:Q(`tr`,{children:[Q(`th`,{children:`Method`}),Q(`th`,{children:`Evidence`}),Q(`th`,{children:`Efficiency`}),Q(`th`,{children:`Convergence`}),Q(`th`,{children:`Best For`})]})}),Q(`tbody`,{children:[Q(`tr`,{children:[Q(`td`,{children:`Forward`}),Q(`td`,{children:`No`}),Q(`td`,{children:`High`}),Q(`td`,{children:`Immediate`}),Q(`td`,{children:`Generating training data, no-evidence estimates`})]}),Q(`tr`,{children:[Q(`td`,{children:`Rejection`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Low (if evidence is rare)`}),Q(`td`,{children:`Immediate`}),Q(`td`,{children:`Simple evidence, high-probability evidence`})]}),Q(`tr`,{children:[Q(`td`,{children:`Likelihood-Weighted`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`Medium`}),Q(`td`,{children:`Immediate`}),Q(`td`,{children:`General evidence, moderate networks`})]}),Q(`tr`,{children:[Q(`td`,{children:`Gibbs`}),Q(`td`,{children:`Yes`}),Q(`td`,{children:`High`}),Q(`td`,{children:`Requires burn-in`}),Q(`td`,{children:`Large networks, complex evidence`})]})]})]}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-3`,children:`Tutorial 3`}),`: Use sampled data for structure learning`]}),Q(`li`,{children:[Q($,{to:`tutorial-7`,children:`Tutorial 7`}),`: Advanced model types`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-7`,children:[Q(`h2`,{children:`Tutorial 7: Advanced Models`}),Q(`p`,{children:`Beyond the standard BayesianNetwork, datascience provides several specialized model types for different use cases. This tutorial covers Dynamic Bayesian Networks, Markov Networks, Naive Bayes classifiers, and Structural Equation Models.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:[Q(`li`,{children:`Completed Tutorials 1-3`}),Q(`li`,{children:`Understanding of the BayesianNetwork model type`})]}),Q(`h3`,{children:`Dynamic Bayesian Networks (DBN)`}),Q(`p`,{children:`A Dynamic Bayesian Network models temporal processes. It defines relationships between variables at time t and time t+1, allowing inference about how state evolves over time.`}),Q(`pre`,{children:Q(`code`,{children:`import (
    "github.com/asymmetric-effort/datascience/lib/pgm/models"
    "github.com/asymmetric-effort/datascience/lib/pgm/inference"
)

// Create a 2-time-slice DBN
dbn := models.NewDynamicBayesianNetwork()

// Add nodes at time 0 and time 1
dbn.AddNode("Weather_0")
dbn.AddNode("Weather_1")
dbn.AddNode("Mood_0")
dbn.AddNode("Mood_1")

// Intra-slice edges (within same time step)
dbn.AddEdge("Weather_0", "Mood_0")
dbn.AddEdge("Weather_1", "Mood_1")

// Inter-slice edges (across time steps)
dbn.AddEdge("Weather_0", "Weather_1")  // weather persists
dbn.AddEdge("Mood_0", "Mood_1")        // mood persists

// Set states and CPDs (similar to regular BN)
// ... (states and CPDs for each node)

// DBN-specific inference
dbnInf, _ := inference.NewDBNInference(dbn)
// Query variables at future time steps
// result, _ := dbnInf.Query(...)
_ = dbnInf`})}),Q(`h3`,{children:`Markov Networks (Undirected Graphical Models)`}),Q(`p`,{children:`Markov Networks use undirected edges and factor potentials instead of directed edges and CPDs. They are useful when the direction of influence is unknown or symmetric.`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/pgm/models"

// Create a Markov Network
mn := models.NewMarkovNetwork()

// Add nodes
mn.AddNode("A")
mn.AddNode("B")
mn.AddNode("C")

// Add undirected edges
mn.AddEdge("A", "B")
mn.AddEdge("B", "C")

// Set states
mn.SetStates("A", []string{"a0", "a1"})
mn.SetStates("B", []string{"b0", "b1"})
mn.SetStates("C", []string{"c0", "c1"})

// Add factor potentials (not CPDs -- these don't need to sum to 1)
// Factor over {A, B}
mn.AddFactor(factors.NewDiscreteFactor(
    []string{"A", "B"},
    []int{2, 2},
    []float64{30, 5, 1, 10},  // compatibility scores
))

// Factor over {B, C}
mn.AddFactor(factors.NewDiscreteFactor(
    []string{"B", "C"},
    []int{2, 2},
    []float64{100, 1, 1, 100},
))

// Run inference (Belief Propagation works on Markov Networks)
// The partition function normalizes the potentials to probabilities`})}),Q(`h3`,{children:`Naive Bayes Classifier`}),Q(`p`,{children:`A Naive Bayes classifier is a Bayesian network where a single class variable is the parent of all feature variables, which are assumed conditionally independent given the class.`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/pgm/models"

// Create a Naive Bayes classifier
nb := models.NewNaiveBayes()

// Set the class variable and features
nb.AddNode("Spam")           // class variable
nb.AddNode("HasLink")        // feature 1
nb.AddNode("HasAttachment")  // feature 2
nb.AddNode("LongSubject")    // feature 3

// Edges from class to all features (automatic in NaiveBayes)
nb.AddEdge("Spam", "HasLink")
nb.AddEdge("Spam", "HasAttachment")
nb.AddEdge("Spam", "LongSubject")

// Set states and CPDs
nb.SetStates("Spam", []string{"no", "yes"})
nb.SetStates("HasLink", []string{"no", "yes"})
nb.SetStates("HasAttachment", []string{"no", "yes"})
nb.SetStates("LongSubject", []string{"no", "yes"})

nb.SetCPD("Spam", factors.NewTabularCPD(
    "Spam", 2, []float64{0.7, 0.3}, nil, nil,
))
nb.SetCPD("HasLink", factors.NewTabularCPD(
    "HasLink", 2,
    []float64{0.9, 0.2, 0.1, 0.8},  // links much more common in spam
    []string{"Spam"}, []int{2},
))
// ... CPDs for HasAttachment, LongSubject

// Classify: P(Spam | HasLink=yes, HasAttachment=no, LongSubject=yes)
facs, _ := nb.ToMarkovFactors()
ve := inference.NewVariableElimination(facs)
result, _ := ve.Query(
    []string{"Spam"},
    map[string]int{"HasLink": 1, "HasAttachment": 0, "LongSubject": 1},
)
fmt.Println("P(Spam | features):", result.Values().Data())`})}),Q(`h3`,{children:`Structural Equation Models (SEM)`}),Q(`p`,{children:`SEMs define variables as explicit equations of their parents plus noise terms. They are widely used in causal inference and econometrics.`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/pgm/models"

// Create a Structural Equation Model
sem := models.NewSEM()

// Add variables and structural equations
sem.AddNode("Education")
sem.AddNode("Income")
sem.AddNode("Health")

// Causal structure
sem.AddEdge("Education", "Income")
sem.AddEdge("Education", "Health")
sem.AddEdge("Income", "Health")

// Define equations:
// Income = beta1 * Education + noise1
// Health = beta2 * Education + beta3 * Income + noise2
// ...

// SEM-specific estimation
// semEst := learning.NewSEMEstimator(sem, data)
// semEst.Estimate()`})}),Q(`h3`,{children:`Linear Gaussian BN`}),Q(`p`,{children:`For continuous variables with linear relationships and Gaussian noise, use LinearGaussianBN. Each variable is defined as a linear combination of its parents plus Gaussian noise.`}),Q(`pre`,{children:Q(`code`,{children:`import "github.com/asymmetric-effort/datascience/lib/pgm/models"

lgbn := models.NewLinearGaussianBN()
lgbn.AddNode("X")
lgbn.AddNode("Y")
lgbn.AddNode("Z")
lgbn.AddEdge("X", "Y")
lgbn.AddEdge("Y", "Z")

// Set linear Gaussian CPDs
// Y = 0.5 * X + N(0, 1)
lgbn.SetCPD("X", factors.NewLinearGaussianCPD(
    "X", 0.0, 1.0, nil, nil,  // mean=0, variance=1, no parents
))
lgbn.SetCPD("Y", factors.NewLinearGaussianCPD(
    "Y", 0.0, 1.0,
    []string{"X"}, []float64{0.5},  // beta=0.5 for parent X
))
lgbn.SetCPD("Z", factors.NewLinearGaussianCPD(
    "Z", 0.0, 1.0,
    []string{"Y"}, []float64{0.8},  // beta=0.8 for parent Y
))`})}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q($,{to:`tutorial-8`,children:`Tutorial 8`}),`: Using the internal libraries (numgo, scigo, graphgo, tabgo) directly`]}),Q(`li`,{children:[Q(T,{to:`/api`,children:`API Reference`}),`: Full type and method documentation`]})]})]}),Q(`section`,{class:`section`,id:`tutorial-8`,children:[Q(`h2`,{children:`Tutorial 8: Using the Internal Libraries`}),Q(`p`,{children:`datascience's internal libraries -- numgo, scigo, graphgo, and tabgo -- are general-purpose and can be used independently of the PGM layer. This tutorial shows how to use each for common tasks.`}),Q(`h3`,{children:`Prerequisites`}),Q(`ul`,{children:Q(`li`,{children:[`datascience installed: `,Q(`code`,{children:`go get github.com/asymmetric-effort/datascience`})]})}),Q(`h3`,{children:`numgo: Matrix and Array Operations`}),Q(`p`,{children:`numgo provides N-dimensional arrays, matrices, and vectors -- similar to numpy in Python.`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "github.com/asymmetric-effort/datascience/lib/numgo"
)

func main() {
    // === Vectors ===
    v1 := numgo.NewVector([]float64{1, 2, 3, 4, 5})
    v2 := numgo.NewVector([]float64{5, 4, 3, 2, 1})

    dot := v1.Dot(v2)
    fmt.Println("Dot product:", dot)  // 35

    norm := v1.Norm()
    fmt.Println("L2 norm:", norm)     // sqrt(55) ~ 7.416

    // Element-wise operations
    sum := v1.Add(v2)
    fmt.Println("v1 + v2:", sum.Data())  // [6 6 6 6 6]

    // === Matrices ===
    m1 := numgo.NewMatrixFromData(2, 3, []float64{
        1, 2, 3,
        4, 5, 6,
    })
    m2 := numgo.NewMatrixFromData(3, 2, []float64{
        7, 8,
        9, 10,
        11, 12,
    })

    product := m1.Multiply(m2)
    fmt.Println("Matrix product shape:", product.Rows(), "x", product.Cols())
    fmt.Println("Matrix product:", product.Data())

    // Transpose
    transposed := m1.Transpose()
    fmt.Println("Transposed shape:", transposed.Rows(), "x", transposed.Cols())

    // Square matrix operations
    sq := numgo.NewMatrixFromData(3, 3, []float64{
        1, 2, 3,
        0, 1, 4,
        5, 6, 0,
    })
    det := sq.Det()
    fmt.Println("Determinant:", det)  // 1

    inv := sq.Inverse()
    fmt.Println("Inverse:", inv.Data())

    // === NDArrays ===
    arr := numgo.NewNDArray([]int{2, 3, 4})  // 2x3x4 array
    arr.Fill(1.0)
    total := arr.Sum()
    fmt.Println("Sum of 2x3x4 ones:", total)  // 24
}`})}),Q(`h3`,{children:`scigo: Statistical Distributions and Optimization`}),Q(`p`,{children:`scigo provides statistical distributions, optimization, and hypothesis testing -- similar to scipy in Python.`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "github.com/asymmetric-effort/datascience/lib/scigo"
)

func main() {
    // === Normal Distribution ===
    normal := scigo.NewNormal(0, 1)  // standard normal
    fmt.Println("PDF at 0:", normal.PDF(0))      // ~0.3989
    fmt.Println("CDF at 1.96:", normal.CDF(1.96)) // ~0.975
    fmt.Println("PPF at 0.975:", normal.PPF(0.975)) // ~1.96

    // Sample from the distribution
    sample := normal.Sample(1000)
    fmt.Println("Sample mean:", mean(sample))  // ~0

    // === Chi-Squared Distribution ===
    // Used in chi-squared independence tests
    chi2 := scigo.NewChiSquared(5)  // 5 degrees of freedom
    fmt.Println("Chi2 CDF at 11.07:", chi2.CDF(11.07))  // ~0.95
    pValue := 1.0 - chi2.CDF(11.07)
    fmt.Println("p-value:", pValue)

    // === Other Distributions ===
    beta := scigo.NewBeta(2, 5)
    fmt.Println("Beta mean:", beta.Mean())  // 2/7 ~ 0.286

    gamma := scigo.NewGamma(2, 1)
    fmt.Println("Gamma mean:", gamma.Mean())  // 2.0

    studentT := scigo.NewStudentT(10)  // 10 degrees of freedom
    fmt.Println("t CDF at 2.228:", studentT.CDF(2.228))

    // === Optimization ===
    // Minimize f(x) = (x-3)^2 on [0, 10]
    result := scigo.Minimize(func(x float64) float64 {
        return (x - 3) * (x - 3)
    }, 0.0, 10.0)
    fmt.Println("Minimum at x =", result)  // ~3.0
}`})}),Q(`h3`,{children:`graphgo: Graph Algorithms`}),Q(`p`,{children:`graphgo provides directed and undirected graphs with a full suite of algorithms -- similar to networkx in Python.`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "github.com/asymmetric-effort/datascience/lib/graphgo"
)

func main() {
    // === Directed Graph ===
    dg := graphgo.NewDiGraph()
    dg.AddNode("A")
    dg.AddNode("B")
    dg.AddNode("C")
    dg.AddNode("D")
    dg.AddEdge("A", "B")
    dg.AddEdge("A", "C")
    dg.AddEdge("B", "D")
    dg.AddEdge("C", "D")

    // Basic queries
    fmt.Println("Nodes:", dg.Nodes())
    fmt.Println("Successors of A:", dg.Successors("A"))      // [B, C]
    fmt.Println("Predecessors of D:", dg.Predecessors("D"))   // [B, C]

    // Topological sort
    sorted := dg.TopologicalSort()
    fmt.Println("Topological order:", sorted)  // [A, B, C, D] or [A, C, B, D]

    // Cycle detection
    fmt.Println("Has cycle:", dg.HasCycle())  // false

    // === Undirected Graph ===
    ug := graphgo.NewGraph()
    ug.AddEdge("X", "Y")
    ug.AddEdge("Y", "Z")
    ug.AddEdge("Z", "X")
    ug.AddEdge("W", "V")

    fmt.Println("Neighbors of Y:", ug.Neighbors("Y"))  // [X, Z]
    fmt.Println("Degree of Y:", ug.Degree("Y"))         // 2

    // Connected components
    components := ug.ConnectedComponents()
    fmt.Println("Components:", components)  // [[X Y Z], [W V]]

    // === PDAG (Partially Directed) ===
    pdag := graphgo.NewPDAG()
    pdag.AddDirectedEdge("A", "B")
    pdag.AddUndirectedEdge("B", "C")
    // PDAGs represent Markov equivalence classes

    // === Graph Algorithms ===
    // Moral graph (for converting BN to undirected)
    moral := dg.MoralGraph()
    fmt.Println("Moral graph edges:", moral.Edges())

    // D-separation (for BN independence queries)
    // separated := dg.DSeparation([]string{"A"}, []string{"D"}, []string{"B"})
}`})}),Q(`h3`,{children:`tabgo: Data Manipulation`}),Q(`p`,{children:`tabgo provides DataFrames and Series for tabular data -- similar to pandas in Python.`}),Q(`pre`,{children:Q(`code`,{children:`package main

import (
    "fmt"
    "github.com/asymmetric-effort/datascience/lib/tabgo"
)

func main() {
    // === Read CSV ===
    df, err := tabgo.ReadCSV("data.csv")
    if err != nil {
        panic(err)
    }
    fmt.Printf("DataFrame: %d rows, %d columns\\n", df.NRows(), len(df.Columns()))
    fmt.Println("Column names:", df.Columns())

    // === Access Columns ===
    col := df.Column("Age")
    fmt.Println("Unique values:", col.Unique())
    fmt.Println("Value counts:", col.ValueCounts())

    // === Filter Rows ===
    adults := df.Filter(func(row map[string]interface{}) bool {
        age, ok := row["Age"].(int)
        return ok && age >= 18
    })
    fmt.Printf("Adults: %d rows\\n", adults.NRows())

    // === GroupBy ===
    // Group by a column and compute aggregates
    // grouped := df.GroupBy("Category")

    // === Create DataFrame from scratch ===
    newDF := tabgo.NewDataFrame(
        map[string][]interface{}{
            "Name": {"Alice", "Bob", "Charlie"},
            "Score": {95, 87, 92},
        },
    )
    fmt.Println("New DataFrame columns:", newDF.Columns())

    // === Write CSV ===
    tabgo.WriteCSV(newDF, "output.csv")
    fmt.Println("Wrote output.csv")

    // === Other I/O ===
    // tabgo.ReadParquet("data.parquet")
    // tabgo.WriteParquet(df, "output.parquet")
    // tabgo.ReadXLSX("data.xlsx")
}`})}),Q(`h3`,{children:`Combining Libraries`}),Q(`p`,{children:`The libraries work together naturally. Here is an example that loads data with tabgo, computes statistics with scigo, and uses numgo for matrix operations:`}),Q(`pre`,{children:Q(`code`,{children:`// Load data
df, _ := tabgo.ReadCSV("measurements.csv")

// Extract a column as a float slice
values := df.Column("Temperature").Float64()

// Compute statistics with scigo
mean := scigo.Mean(values)
std := scigo.Std(values)
fmt.Printf("Temperature: mean=%.2f, std=%.2f\\n", mean, std)

// Build a correlation matrix with numgo
cols := []string{"Temperature", "Humidity", "Pressure"}
n := len(cols)
corrMatrix := numgo.NewMatrix(n, n)
for i := 0; i < n; i++ {
    for j := 0; j < n; j++ {
        vi := df.Column(cols[i]).Float64()
        vj := df.Column(cols[j]).Float64()
        corrMatrix.Set(i, j, scigo.PearsonCorrelation(vi, vj))
    }
}
fmt.Println("Correlation matrix:", corrMatrix.Data())`})}),Q(`h3`,{children:`What's Next`}),Q(`ul`,{children:[Q(`li`,{children:[Q(T,{to:`/docs`,children:`Documentation`}),`: Full package reference for all types and methods`]}),Q(`li`,{children:[Q(T,{to:`/api`,children:`API Reference`}),`: Detailed API documentation with signatures and examples`]})]})]})]})}function $n(){return Q(Fe,{children:Q(`div`,{id:`app`,children:[Q(Kn,{}),Q(`main`,{children:[Q(w,{path:`/`,component:Jn,exact:!0}),Q(w,{path:`/docs`,component:Yn}),Q(w,{path:`/libraries`,component:Xn}),Q(w,{path:`/api`,component:Zn}),Q(w,{path:`/tutorials`,component:Qn})]}),Q(qn,{})]})})}var er=document.getElementById(`root`);er&&Gn(er).render(Q($n,{}));