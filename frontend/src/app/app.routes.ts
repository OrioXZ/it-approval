import { Routes } from '@angular/router';

import { It03Component } from './pages/it03/it03';

export const routes: Routes = [
  { path: '', redirectTo: 'it03', pathMatch: 'full' },
  { path: 'it03', component: It03Component },
];
