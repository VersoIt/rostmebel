/// <reference types="vite/client" />

declare const __PUBLIC_SITE_URL__: string;

declare module "*.vue" {
  import type { DefineComponent } from "vue";
  const component: DefineComponent<{}, {}, any>;
  export default component;
}
