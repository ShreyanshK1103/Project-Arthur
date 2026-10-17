/**
 * Public site URLs for the dashboard.
 *
 * The API already returns `slug` (projects) and `id` (deployments).
 * Hosted sites are served as subdomains of SITE_BASE_DOMAIN:
 *   {project-slug}.localhost:8080
 *   {first-12-chars-of-deployment-id}.localhost:8080
 */

const SITE_SCHEME =
  (import.meta.env.VITE_SITE_SCHEME as string | undefined) || "http";
const SITE_BASE_DOMAIN =
  (import.meta.env.VITE_SITE_BASE_DOMAIN as string | undefined) ||
  "localhost:8080";

/** Backend uses uuid.String()[:12] as the deployment subdomain. */
export function deploymentSlug(deploymentId: string): string {
  return deploymentId.slice(0, 12);
}

export function siteHost(slug: string): string {
  return `${slug}.${SITE_BASE_DOMAIN}`;
}

export function siteHref(slug: string): string {
  return `${SITE_SCHEME}://${siteHost(slug)}`;
}

export function projectSite(slug: string) {
  return { slug, host: siteHost(slug), href: siteHref(slug) };
}

export function deploymentSite(deploymentId: string) {
  const slug = deploymentSlug(deploymentId);
  return { slug, host: siteHost(slug), href: siteHref(slug) };
}
