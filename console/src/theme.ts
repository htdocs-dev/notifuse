import { theme as antdTheme, type ThemeConfig } from "antd";
import type { ResolvedTheme } from "./themeMode";

/**
 * Ant Design theme for the console (fork patch 6).
 *
 * Surfaces, borders and text point at the CSS variables declared in index.css
 * (`--nf-*`), so one place decides the palette for antd components, Tailwind
 * utilities and inline styles alike. Only the seed colours antd derives from
 * (primary, background and text base) are literals, because the algorithm
 * needs real colours to compute hover, active and fill shades.
 */
export function buildAntdTheme(resolved: ResolvedTheme): ThemeConfig {
  const dark = resolved === "dark";
  const primary = dark ? "#3ECF8E" : "#12A66F";
  return {
    algorithm: dark ? antdTheme.darkAlgorithm : antdTheme.defaultAlgorithm,
    token: {
      colorPrimary: primary,
      colorLink: primary,
      colorInfo: primary,
      colorBgBase: dark ? "#0B0D11" : "#FFFFFF",
      colorTextBase: dark ? "#E6EAF0" : "#101828",
      colorBgLayout: "var(--nf-page)",
      colorBgContainer: "var(--nf-surface)",
      colorBgElevated: "var(--nf-elevated)",
      colorBorder: "var(--nf-border-strong)",
      colorBorderSecondary: "var(--nf-border)",
      colorText: "var(--nf-text)",
      colorTextSecondary: "var(--nf-text-2)",
      colorTextTertiary: "var(--nf-text-3)",
      colorTextDescription: "var(--nf-text-2)",
      fontFamily:
        '"Inter Variable", Inter, system-ui, -apple-system, "Segoe UI", Helvetica, Arial, sans-serif',
      fontSize: 13,
      borderRadius: 8,
      borderRadiusLG: 12,
      borderRadiusSM: 6,
      controlHeight: 34,
      controlOutline: "var(--nf-ring)",
      controlOutlineWidth: 3,
      boxShadowSecondary: dark
        ? "0 8px 24px rgba(0, 0, 0, 0.45), 0 0 0 1px var(--nf-border)"
        : "0 8px 24px rgba(16, 24, 40, 0.1), 0 0 0 1px var(--nf-border)",
    },
    components: {
      Layout: {
        bodyBg: "transparent",
        headerBg: "var(--nf-surface)",
        lightSiderBg: "var(--nf-surface)",
        siderBg: "var(--nf-surface)",
      },
      Menu: {
        itemBg: "transparent",
        subMenuItemBg: "transparent",
        itemColor: "var(--nf-text-2)",
        itemHoverBg: "var(--nf-hover)",
        itemHoverColor: "var(--nf-text)",
        itemSelectedBg: "var(--nf-accent-soft)",
        itemSelectedColor: primary,
        itemBorderRadius: 8,
        itemHeight: 36,
        itemMarginInline: 8,
        itemMarginBlock: 2,
        iconSize: 14,
        activeBarBorderWidth: 0,
      },
      Button: {
        fontWeight: 500,
        // Mint carries dark text; white on it would fail contrast.
        primaryColor: dark ? "#07130D" : "#FFFFFF",
        primaryShadow: "none",
        defaultShadow: "none",
        dangerShadow: "none",
        defaultBg: "var(--nf-surface)",
        defaultBorderColor: "var(--nf-border-strong)",
        defaultHoverBg: "var(--nf-hover)",
      },
      Card: {
        headerFontSize: 15,
        // Card sizes every corner from borderRadiusLG alone. The other radius tokens never
        // reached the card, and overriding a global token per component now emits a scoped
        // CSS variable, so they would only reround the buttons, inputs and popups nested
        // inside cards.
        borderRadiusLG: 12,
        colorBorderSecondary: "var(--nf-border)",
        colorBgContainer: "var(--nf-surface)",
      },
      Table: {
        headerBg: "transparent",
        // Sized and coloured through the cell/header tokens instead of the global fontSize
        // and colorTextHeading: those cascade out of the table wrapper as CSS variables and
        // would shrink and recolour every antd component rendered inside a cell.
        cellFontSize: 13,
        cellFontSizeMD: 13,
        cellFontSizeSM: 12,
        cellPaddingBlock: 12,
        headerColor: "var(--nf-text-3)",
        footerColor: "var(--nf-text-2)",
        borderColor: "var(--nf-border)",
        colorBgContainer: "transparent",
        rowHoverBg: "var(--nf-hover)",
        rowSelectedBg: "var(--nf-accent-soft)",
        rowSelectedHoverBg: "var(--nf-accent-soft)",
        // The container is transparent, so antd's default sort-highlight fills resolve to
        // opaque black on the sorted column and hovered sortable headers. Keep them
        // transparent to match the flat table style; the sort arrow still signals order.
        headerSortActiveBg: "transparent",
        headerSortHoverBg: "transparent",
        bodySortBg: "transparent",
      },
      Input: {
        colorBgContainer: "var(--nf-input-bg)",
        activeShadow: "0 0 0 3px var(--nf-ring)",
      },
      InputNumber: {
        colorBgContainer: "var(--nf-input-bg)",
        activeShadow: "0 0 0 3px var(--nf-ring)",
      },
      Select: {
        colorBgContainer: "var(--nf-input-bg)",
        optionSelectedBg: "var(--nf-accent-soft)",
      },
      DatePicker: {
        colorBgContainer: "var(--nf-input-bg)",
        activeShadow: "0 0 0 3px var(--nf-ring)",
      },
      Tabs: {
        inkBarColor: primary,
        itemColor: "var(--nf-text-2)",
        itemSelectedColor: "var(--nf-text)",
        itemHoverColor: "var(--nf-text)",
      },
      Tag: {
        defaultBg: "var(--nf-chip)",
        defaultColor: "var(--nf-text-2)",
        borderRadiusSM: 6,
      },
      Drawer: {
        // Drawer paints its panel straight from colorBgElevated and exposes no background
        // token of its own, so this override has to stay on the global token.
        colorBgElevated: "var(--nf-elevated)",
      },
      Modal: {
        contentBg: "var(--nf-elevated)",
        headerBg: "var(--nf-elevated)",
      },
      Timeline: {
        dotBg: "var(--nf-surface)",
      },
      Tooltip: {
        colorBgSpotlight: dark ? "#262C36" : "#101828",
      },
    },
  };
}
