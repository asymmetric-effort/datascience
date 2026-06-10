import { createElement, Link, useRouter } from "@asymmetric-effort/specifyjs";

export function Nav() {
  const { pathname } = useRouter();

  return (
    <nav class="nav">
      <div class="nav-inner">
        <Link to="/" class="nav-brand">
          <img src="/docs/img/logo.png" alt="datascience logo" width="28" height="28" />
          datascience
        </Link>
        <div class="nav-links">
          <Link to="/" class={pathname === "/" ? "active" : ""}>Home</Link>
          <Link to="/docs" class={pathname === "/docs" ? "active" : ""}>Docs</Link>
          <Link to="/tutorials" class={pathname === "/tutorials" ? "active" : ""}>Tutorials</Link>
          <Link to="/api" class={pathname === "/api" ? "active" : ""}>API</Link>
          <Link to="/libraries" class={pathname === "/libraries" ? "active" : ""}>Libraries</Link>
        </div>
      </div>
    </nav>
  );
}
