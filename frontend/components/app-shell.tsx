"use client"

import { Outlet } from "react-router-dom"
import { MonitorHeader } from "@/components/monitor/monitor-header"
import { DockBar } from "@/components/monitor/dock-bar"

/**
 * AppShell 是所有路由共享的外壳：顶部 header + 中间 Outlet（+ 可选底部 dock）。
 *
 * 底部 Dock 承载二级页面入口和全局“新增渠道”动作。
 * 这些入口在 header 中没有重复出现，隐藏 Dock 会让设置/通知/打码页不可发现。
 */
const SHOW_DOCK = false

export function AppShell() {
  return (
    <div className="min-h-screen bg-background">
      <MonitorHeader />
      <main
        className={
          SHOW_DOCK
            ? "mx-auto max-w-360 space-y-5 px-5 py-5 pb-24"
            : "mx-auto max-w-360 space-y-5 px-5 py-5"
        }
      >
        <Outlet />
      </main>
      {SHOW_DOCK ? <DockBar /> : null}
    </div>
  )
}
