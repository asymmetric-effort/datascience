import { createElement, Link, useRouter } from "@asymmetric-effort/specifyjs";

export function Nav() {
  const { pathname } = useRouter();

  return (
    <nav class="nav">
      <div class="nav-inner">
        <Link href="/" class="nav-brand">pgmgo</Link>
        <div class="nav-links">
          <Link href="/" class={pathname === "/" ? "active" : ""}>Home</Link>
          <Link href="/docs" class={pathname === "/docs" ? "active" : ""}>Docs</Link>
          <Link href="/cli" class={pathname === "/cli" ? "active" : ""}>CLI</Link>
          <Link href="/api" class={pathname === "/api" ? "active" : ""}>API</Link>
        </div>
      </div>
    </nav>
  );
}
