<script setup lang="ts">
 import { ref, onMounted } from 'vue'
 import { Viewer, Editor } from '@bytemd/vue-next'
 // import Ws from './components/ws.vue'
 import breaks from '@bytemd/plugin-breaks'
 import frontmatter from '@bytemd/plugin-frontmatter'
 import highlight from '@bytemd/plugin-highlight'
 import gfm from '@bytemd/plugin-gfm'
 import mediumZoom from '@bytemd/plugin-medium-zoom'
 import mermaid from '@bytemd/plugin-mermaid'
 import 'bytemd/dist/index.css'
 import 'github-markdown-css'
 import 'highlight.js/styles/vs.css'
 import 'katex/dist/katex.css'
 import {
   ArrayQueue,
   ConstantBackoff,
   Websocket,
   WebsocketBuilder,
   WebsocketEvent,
 } from "websocket-ts"


 let content = ref('')

 let enabled = {
   breaks: false,
   frontmatter: true,
   gemoji: true,
   gfm: true,
   highlight: true,
   math: true,
   'medium-zoom': true,
   mermaid: true,
 }

 const plugins = [
   enabled.breaks && breaks(),
   enabled.frontmatter && frontmatter(),
   enabled.highlight && highlight(),
   gfm(),
 ]

 const handleChange = (v: string) => {
   content.value = v
 }


 onMounted(() => {
   let urlParams = new URLSearchParams(window.location.search);
   if (!urlParams.get('server')) {
     return
   }
   if (!urlParams.get('path')) {
     return
   }

   const host = urlParams.get('server') //"127.0.0.1:8080"
   const path = urlParams.get('path')

   const ws = new WebsocketBuilder("ws://" + host + "/stream")
   .withBuffer(new ArrayQueue())           // buffer messages when disconnected
   .withBackoff(new ConstantBackoff(10000)) // retry every 1s
   .build();

   const echoOnMessage = (i: Websocket, ev: MessageEvent) => {
     content.value = ev.data
   };

   // Add event listeners
   ws.addEventListener(WebsocketEvent.open, (i: Websocket, ev: MessageEvent) => {
     console.log("opened!")
     i.send(path)
   });
   ws.addEventListener(WebsocketEvent.close, () => console.log("closed!"));
   ws.addEventListener(WebsocketEvent.message, echoOnMessage);
 });


</script>

<template>
  <main>
    <Viewer :value="content" :plugins="plugins" @change="handleChange" />
  </main>
</template>

<style>
 #app {
   max-width: 1280px;
   margin: 0 auto;
   padding: 2rem;
   font-weight: normal;
 }
</style>
