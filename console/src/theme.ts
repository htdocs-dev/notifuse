import { theme as antdTheme, type ThemeConfig } from "antd";
import type { ResolvedTheme } from "./themeMode";

/**
 * Ant Design theme for the console.
 *
 * Surfaces point at the CSS variables declared in index.css (`--nf-surface`,
 * `--nf-border`, Tailwind's `--color-*`), so one place decides the palette for
 * antd components, Tailwind utilities and inline styles alike. Only the seed
 * colours antd derives from (primary, background base) are literals, because
 * the algorithm needs real colours to compute hover and active shades.
 */
export function buildAntdTheme(resolved: ResolvedTheme): ThemeConfig {
  const dark = resolved === "dark";
  return {
    algorithm: dark ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
    token: {
      // Teal accent on cool greys (fork patch 6). Lighter in dark mode so the
      // hover and active shades antd derives keep contrast on the navy ground.
      colorPrimary: dark ? "#2DD4BF" : "#0D9488",
      colorLink: dark ? "#2DD4BF" : "#0D9488",
      ...(dark ? { colorBgBase: "#0F141A", colorTextBase: "#E6EBF0" } : {}),
    },
    components: {
      Layout: {
        bodyBg: "var(--nf-surface)",
        lightSiderBg: "var(--nf-surface)",
        siderBg: "var(--nf-surface)",
      },
      Card: {
        headerFontSize: 16,
        // Card sizes every corner from borderRadiusLG alone. The other radius tokens never
        // reached the card, and overriding a global token per component now emits a scoped
        // CSS variable, so they would only reround the buttons, inputs and popups nested
        // inside cards.
        borderRadiusLG: 4,
        colorBorderSecondary: "var(--color-gray-200)",
        colorBgContainer: "var(--nf-surface)",
      },
      Table: {
        headerBg: "transparent",
        // Sized and coloured through the cell/header tokens instead of the global fontSize
        // and colorTextHeading: those cascade out of the table wrapper as CSS variables and
        // would shrink and recolour every antd component rendered inside a cell.
        cellFontSize: 12,
        cellFontSizeMD: 12,
        cellFontSizeSM: 12,
        headerColor: "var(--color-slate-700)",
        footerColor: "var(--color-slate-700)",
        colorBgContainer: "transparent",
        rowHoverBg: "transparent",
        // The container is transparent, so antd's default sort-highlight fills resolve to
        // opaque black on the sorted column and hovered sortable headers. Keep them
        // transparent to match the flat table style; the sort arrow still signals order.
        headerSortActiveBg: "transparent",
        headerSortHoverBg: "transparent",
        bodySortBg: "transparent",
      },
      Drawer: {
        // Drawer paints its panel straight from colorBgElevated and exposes no background
        // token of its own, so this override has to stay on the global token.
        colorBgElevated: "var(--nf-surface)",
      },
      Modal: {
        // contentBg is the only thing Modal derived from colorBgElevated, and setting it
        // directly keeps the override from cascading into the dialog's popups and sliders.
        contentBg: "var(--nf-surface)",
      },
      Timeline: {
        dotBg: "var(--nf-surface)",
      },
    },
  };
}
