package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

const (
	compactWidth = 1394
	monthWidth   = 49.66
)

func (a *app) dashboardPage(c *ui.Context) {
	p := a.colors()
	width, _ := c.Size()
	columns := 4
	if width < compactWidth {
		columns = 2
	}
	ui.Grid(c).Columns(columns).GapX(17).GapY(17).Children(func() {
		for _, m := range metrics {
			metricCard(c, p, a.metricView(m))
		}
	})
	if width < compactWidth {
		ui.Column(c).Gap(18).Children(func() {
			a.revenueCard(c).Height(289)
			a.categoryCard(c).Height(289)
		})
	} else {
		ui.Row(c).Gap(18).Children(func() {
			a.revenueCard(c).Height(289).Grow(1).Basis(0).MinWidth(0)
			a.categoryCard(c).Height(289).Width(403)
		})
	}
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Children(func() {
			cardTitle(c, "Recent Orders")
			view := moreButton(c, p, "View All", arrowIcon)
			if view.Clicked() {
				a.page, a.query = pageOrders, ""
			}
		})
		recent := filterOrders(a.query)
		if len(recent) > 5 {
			recent = recent[:5]
		}
		a.ordersTable(c, recent)
	})
}

func (a *app) metricView(m metric) metricView {
	value, previous := m.at(a.month)
	trend := m.trend()
	return metricView{
		label:    m.label,
		icon:     m.icon,
		value:    m.format(value),
		against:  "vs. last month",
		delta:    (value - previous) / previous * 100,
		up:       arrowUpIcon,
		down:     arrowDownIcon,
		trend:    trend,
		describe: func(i int) string { return trendMonths[i] + ": " + m.format(trend[i]) },
	}
}

func (a *app) revenueCard(c *ui.Context) ui.Element {
	p := a.colors()
	return card(c, p).Padding(15, 19, 17, 22).Gap(26).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Children(func() {
			cardTitle(c, "Revenue (Monthly)")
			segmented(c, p, "month", "Month", &a.month, months, monthWidth)
		})
		areaChart(c, p, areaChartView{
			label:     "Monthly revenue",
			names:     months,
			values:    metrics[0].monthly[:],
			ceiling:   150000,
			step:      25000,
			axisLabel: axisDollars,
			shaded:    a.month,
			describe:  func(i int) string { return months[i] + ": " + thousands(metrics[0].monthly[i]) },
		})
	})
}

func (a *app) categoryCard(c *ui.Context) ui.Element {
	p := a.colors()
	shares := categoryShares[a.month]
	total, _ := metrics[0].at(a.month)
	return card(c, p).Padding(16, 19, 0, 20).Gap(23).Children(func() {
		cardTitle(c, "Sales by Category")
		donut(c, p, donutView{
			label:    "Sales by category",
			names:    categories,
			shares:   shares,
			center:   thousands(total),
			caption:  "Total",
			describe: func(i int) string { return categories[i] + ": " + thousands(total*shares[i]/100) },
		})
	})
}

func (a *app) ordersTable(c *ui.Context, rows []order) {
	p := a.colors()
	columns := []column{{"Order ID", 203}, {"Customer", 254}, {"Status", 289}, {"Amount", 163}, {"Date", 192}}
	table(c, p, "Orders", "Nothing matches your search.", columns, len(rows), func(row, col int) {
		o := rows[row]
		switch col {
		case 0:
			open := textLink(c, p, o.id)
			if open.Clicked() {
				a.page, a.query = pageOrders, o.id
			}
		case 1:
			cellText(c, o.customer)
		case 2:
			statusPill(c, p.status(o.status), o.status)
		case 3:
			cellText(c, cents(o.amount))
		case 4:
			cellText(c, o.date)
		}
	})
}

func (p *palette) status(status string) tone {
	switch status {
	case statusCompleted:
		return p.good
	case statusCancelled:
		return p.bad
	}
	return p.info
}

func (a *app) analyticsPage(c *ui.Context) {
	p := a.colors()
	a.revenueCard(c).Height(360)
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		cardTitle(c, "Monthly Breakdown")
		columns := []column{{"Month", 1}, {"Total Revenue", 1}, {"Active Users", 1}, {"Orders", 1}, {"Conversion Rate", 1}}
		table(c, p, "Monthly breakdown", "Nothing matches your search.", columns, len(months), func(row, col int) {
			if col == 0 {
				cellText(c, months[row])
				return
			}
			m := metrics[col-1]
			cellText(c, m.format(m.monthly[row]))
		})
	})
}

func (a *app) ordersPage(c *ui.Context) {
	p := a.colors()
	rows := filterOrders(a.query)
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		ui.Row(c).Justify(ui.SpaceBetween).Height(30).Children(func() {
			cardTitle(c, "All Orders")
			ui.Text(c, fmt.Sprintf("%d of %d", len(rows), len(orders))).FontSize(14).TextColor(p.muted)
		})
		a.ordersTable(c, rows)
	})
}

func (a *app) customersPage(c *ui.Context) {
	p := a.colors()
	rows := summarizeCustomers(filterOrders(a.query))
	card(c, p).Padding(15, 16, 20, 22).Gap(13).Children(func() {
		ui.Row(c).Height(30).Children(func() {
			cardTitle(c, "All Customers")
		})
		columns := []column{{"Customer", 2}, {"Orders", 1}, {"Total Spent", 1}, {"Last Order", 1}}
		table(c, p, "Customers", "Nothing matches your search.", columns, len(rows), func(row, col int) {
			customer := rows[row]
			switch col {
			case 0:
				cellText(c, customer.name)
			case 1:
				cellText(c, fmt.Sprint(customer.orders))
			case 2:
				cellText(c, cents(customer.spent))
			case 3:
				cellText(c, customer.last)
			}
		})
	})
}

func (a *app) settingsPage(c *ui.Context) {
	p := a.colors()
	card(c, p).Padding(15, 21, 20, 22).Gap(13).Children(func() {
		cardTitle(c, "Appearance")
		ui.Row(c).Height(38).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Dark mode").FontSize(14.5).TextColor(p.body)
			themeSwitch(c, p, &a.dark, "Use dark mode", sunIcon, moonIcon, a.applyTheme)
		})
	})
}
