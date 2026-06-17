"use client"

import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from "recharts"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { useDashboardSummary } from "@/lib/queries"
import { yuan } from "@/lib/format"
import { cn } from "@/lib/utils"

const COLORS = [
  "var(--brand)",
  "var(--success)",
  "var(--warning)",
  "var(--danger)",
  "var(--chart-5)",
  "var(--muted-foreground)",
]

interface PieItem {
  id: number
  name: string
  value: number
  color: string
  pct: number
  failed: boolean
}

interface TooltipPayloadItem {
  payload: PieItem
}

function BalanceTooltip({
  active,
  payload,
}: {
  active?: boolean
  payload?: TooltipPayloadItem[]
}) {
  if (!active || !payload?.length) return null
  const item = payload[0].payload
  return (
    <div className="rounded-lg border border-border bg-popover px-3 py-2 shadow-md">
      <p className="text-xs font-medium text-foreground">{item.name}</p>
      <p className="mt-0.5 text-sm font-semibold tabular-nums text-foreground">
        {yuan(item.value)}
      </p>
      <p className="mt-0.5 text-xs text-muted-foreground">{item.pct.toFixed(1)}%</p>
    </div>
  )
}

export function BalanceOverview() {
  const summary = useDashboardSummary()
  const channels = summary.data?.channels ?? []
  const total = summary.data?.total_balance ?? 0
  const items: PieItem[] = channels
    .filter((c) => c.last_balance != null && Number.isFinite(c.last_balance) && c.last_balance > 0)
    .map((c, index) => ({
      id: c.id,
      name: c.name,
      value: c.last_balance ?? 0,
      color: COLORS[index % COLORS.length],
      pct: total > 0 ? ((c.last_balance ?? 0) / total) * 100 : 0,
      failed: !!c.last_error,
    }))
    .sort((a, b) => b.value - a.value)

  return (
    <Card className="border border-border shadow-none lg:h-100">
      <CardHeader className="flex shrink-0 flex-row items-center justify-between pb-2">
        <CardTitle className="text-base font-semibold">{"余额概览"}</CardTitle>
        <span className="text-xs text-muted-foreground">{yuan(total)}</span>
      </CardHeader>
      <CardContent className="flex min-h-0 flex-1 flex-col">
        {summary.loading ? (
          <div className="flex h-full items-center justify-center text-xs text-muted-foreground">
            {"加载中…"}
          </div>
        ) : items.length === 0 ? (
          <div className="flex h-full items-center justify-center text-xs text-muted-foreground">
            {"暂无余额采样，等待下次扫描或手动刷新"}
          </div>
        ) : (
          <>
            <div className="relative min-h-0 w-full flex-1">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart margin={{ top: 4, right: 8, bottom: 4, left: 8 }}>
                  <Pie
                    data={items}
                    dataKey="value"
                    nameKey="name"
                    innerRadius="56%"
                    outerRadius="82%"
                    paddingAngle={2}
                    stroke="var(--background)"
                    strokeWidth={3}
                  >
                    {items.map((item) => (
                      <Cell key={item.id} fill={item.color} />
                    ))}
                  </Pie>
                  <Tooltip content={<BalanceTooltip />} />
                </PieChart>
              </ResponsiveContainer>
              <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
                <div className="text-center">
                  <p className="text-[11px] text-muted-foreground">{"总余额"}</p>
                  <p className="mt-0.5 text-lg font-semibold tabular-nums text-foreground">
                    {yuan(total)}
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-3 flex shrink-0 flex-wrap items-center gap-x-5 gap-y-2 border-t border-border pt-3">
              {items.map((c) => (
                <span key={c.id} className="inline-flex items-center gap-1.5 text-xs">
                  <span
                    className={cn(
                      "size-2 rounded-full",
                      c.failed ? "bg-danger" : "",
                    )}
                    style={c.failed ? undefined : { backgroundColor: c.color }}
                  />
                  <span className="font-medium text-foreground">{c.name}</span>
                  <span className="tabular-nums text-muted-foreground">{yuan(c.value)}</span>
                </span>
              ))}
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}
