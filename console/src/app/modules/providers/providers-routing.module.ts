import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ProviderType } from 'src/app/proto/generated/zitadel/idp_pb';
import { ProviderAppleComponent } from './provider-apple/provider-apple.component';
import { ProviderAzureADComponent } from './provider-azure-ad/provider-azure-ad.component';
import { ProviderGithubESComponent } from './provider-github-es/provider-github-es.component';
import { ProviderGithubComponent } from './provider-github/provider-github.component';
import { ProviderGitlabSelfHostedComponent } from './provider-gitlab-self-hosted/provider-gitlab-self-hosted.component';
import { ProviderGitlabComponent } from './provider-gitlab/provider-gitlab.component';
import { ProviderGoogleComponent } from './provider-google/provider-google.component';
import { ProviderJWTComponent } from './provider-jwt/provider-jwt.component';
import { ProviderLDAPComponent } from './provider-ldap/provider-ldap.component';
import { ProviderOAuthComponent } from './provider-oauth/provider-oauth.component';
import { ProviderOIDCComponent } from './provider-oidc/provider-oidc.component';
import { ProviderSamlSpComponent } from './provider-saml-sp/provider-saml-sp.component';

const typeMap = {
  [ProviderType.PROVIDER_TYPE_AZURE_AD]: { path: 'azure-ad', component: ProviderAzureADComponent },
  [ProviderType.PROVIDER_TYPE_GITHUB]: { path: 'github', component: ProviderGithubComponent },
  [ProviderType.PROVIDER_TYPE_GITHUB_ES]: { path: 'github-es', component: ProviderGithubESComponent },
  [ProviderType.PROVIDER_TYPE_GITLAB]: { path: 'gitlab', component: ProviderGitlabComponent },
  [ProviderType.PROVIDER_TYPE_GITLAB_SELF_HOSTED]: {
    path: 'gitlab-self-hosted',
    component: ProviderGitlabSelfHostedComponent,
  },
  [ProviderType.PROVIDER_TYPE_GOOGLE]: { path: 'google', component: ProviderGoogleComponent },
  [ProviderType.PROVIDER_TYPE_JWT]: { path: 'jwt', component: ProviderJWTComponent },
  [ProviderType.PROVIDER_TYPE_OAUTH]: { path: 'oauth', component: ProviderOAuthComponent },
  [ProviderType.PROVIDER_TYPE_OIDC]: { path: 'oidc', component: ProviderOIDCComponent },
  [ProviderType.PROVIDER_TYPE_LDAP]: { path: 'ldap', component: ProviderLDAPComponent },
  [ProviderType.PROVIDER_TYPE_APPLE]: { path: 'apple', component: ProviderAppleComponent },
  [ProviderType.PROVIDER_TYPE_SAML]: { path: 'saml', component: ProviderSamlSpComponent },
};

const routes: Routes = Object.values(typeMap).map(({ path, component }) => {
  return {
    path,
    children: [
      {
        path: 'create',
        component,
      },
      {
        path: ':id',
        component,
      },
    ],
  };
});

// Discord, Twitch and Kick don't have dedicated provider types (that would require proto changes
// and regenerating stubs across the Go backend, this console and the login app). Instead they're
// presets on top of the existing Generic OAuth2 (Discord, Kick) and Generic OIDC (Twitch) provider
// types: same create/edit forms and storage, just pre-filled with the right endpoints/scopes so
// nobody has to look those up by hand. See ProviderOAuthComponent/ProviderOIDCComponent presets.
const presetRoutes: Routes = [
  { path: 'discord', preset: 'discord', component: ProviderOAuthComponent },
  { path: 'kick', preset: 'kick', component: ProviderOAuthComponent },
  { path: 'twitch', preset: 'twitch', component: ProviderOIDCComponent },
].map(({ path, preset, component }) => {
  return {
    path,
    children: [
      {
        path: 'create',
        component,
        data: { preset },
      },
      {
        path: ':id',
        component,
      },
    ],
  };
});

@NgModule({
  imports: [RouterModule.forChild([...routes, ...presetRoutes])],
  exports: [RouterModule],
})
export class ProvidersRoutingModule {}
