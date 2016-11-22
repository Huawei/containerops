/*
 * Angular bootstraping
 */
// import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
// import { decorateModuleRef } from './app/environment';
// import { bootloader } from '@angularclass/hmr';
import {platformBrowserDynamic} from '@angular/platform-browser-dynamic';
/*
 * App Module
 * our top level module that holds all of our components
 */
import { AppModule } from './app/app.module';

/*
 * Bootstrap our Angular app with a top level NgModule
 */
// export function main(): Promise<any> {
//   return platformBrowserDynamic()
//     .bootstrapModule(App)
//     // .then(decorateModuleRef)
//     .catch(err => console.error(err));

// }


// bootloader(main);
// main();
document.addEventListener('DOMContentLoaded', function main(): void {
  platformBrowserDynamic().bootstrapModule(AppModule)
  .catch(err => console.error(err));
});
