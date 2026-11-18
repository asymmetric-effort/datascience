import { createElement, Link, useRouter, useState } from "@asymmetric-effort/specifyjs";

export function Nav() {
  const { pathname } = useRouter();
  const [menuOpen, setMenuOpen] = useState(false);

  const toggleMenu = () => setMenuOpen(!menuOpen);
  const closeMenu = () => setMenuOpen(false);

  return (
    <nav class="nav">
      <div class="nav-inner">
        <Link to="/" class="nav-brand" onClick={closeMenu}>
          <img src="/docs/img/logo.png" alt="datascience.asymmetric-effort.com logo" width="28" height="28" />
          datascience.asymmetric-effort.com
        </Link>
        <button class="nav-toggle" onClick={toggleMenu} aria-label="Toggle navigation menu" aria-expanded={menuOpen ? "true" : "false"}>
          {menuOpen ? "\u2715" : "\u2630"}
        </button>
        <div class={menuOpen ? "nav-links open" : "nav-links"}>
          <Link to="/" class={pathname === "/" ? "active" : ""} onClick={closeMenu}>Home</Link>
          <Link to="/docs" class={pathname === "/docs" ? "active" : ""} onClick={closeMenu}>Docs</Link>
          <Link to="/tutorials" class={pathname === "/tutorials" ? "active" : ""} onClick={closeMenu}>Tutorials</Link>
          <Link to="/api" class={pathname === "/api" ? "active" : ""} onClick={closeMenu}>API</Link>
          <Link to="/libraries" class={pathname === "/libraries" ? "active" : ""} onClick={closeMenu}>Libraries</Link>
        </div>
      </div>
    </nav>
  );
}
