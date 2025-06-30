import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

import { StatsCard } from '@/components/dashboard/StatsCard'
import { AppointmentCard } from '@/components/dashboard/AppointmentCard'
import { BottomNavigation } from '@/components/dashboard/BottomNavigation'
import { ActiveSessionCard } from '@/components/dashboard/ActiveSessionCard'
import {
  calculateScheduleStats,
  formatDuration,
  getActiveSchedule,
} from '@/shared/dashboard'
import { useInfiniteSchedules } from '@/hooks/useSchedules'
import { useCallback, useEffect, useRef, useState } from 'react'

function DesktopHeader() {
  return (
    <div className="hidden md:flex bg-[#B8E5E1] px-8 py-4 justify-between items-center">
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-[#0D5D59] rounded flex items-center justify-center">
            <span className="text-white font-bold text-sm">🏠</span>
          </div>
          <span className="font-semibold text-gray-800">Careviah</span>
        </div>
      </div>
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-gray-400 rounded-full"></div>
          <div>
            <div className="text-sm font-medium text-gray-800">Admin A</div>
            <div className="text-xs text-gray-600">admin@careviah.co.in</div>
          </div>
        </div>
      </div>
    </div>
  )
}

function Dashboard() {
  const {
    data,
    error,
    isLoading,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteSchedules(5) // Use the new infinite query hook with a page size

  const [duration, setDuration] = useState<string>('')

  // Flatten the data from all pages into a single array for rendering
  const allSchedules =
    data?.pages.flatMap((page) => (page.success ? page.data || [] : [])) || []

  const stats = calculateScheduleStats(allSchedules) // Calculate stats from all loaded schedules

  const activeSchedule = getActiveSchedule(allSchedules) // Get active schedule from all loaded schedules

  useEffect(() => {
    if (activeSchedule && activeSchedule.start_time) {
      const interval = setInterval(() => {
        setDuration(formatDuration(String(activeSchedule.start_time)))
      }, 1000)

      return () => clearInterval(interval)
    }
  }, [activeSchedule])

  // Infinite scroll logic
  const observerElem = useRef<HTMLDivElement>(null)

  const handleObserver = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      const target = entries[0]
      if (target.isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage()
      }
    },
    [fetchNextPage, hasNextPage, isFetchingNextPage],
  )

  useEffect(() => {
    const observer = new IntersectionObserver(handleObserver, {
      root: null, // viewport
      rootMargin: '0px',
      threshold: 0.1, // Trigger when 10% of the target is visible
    })

    if (observerElem.current) {
      observer.observe(observerElem.current)
    }

    return () => {
      if (observerElem.current) {
        observer.unobserve(observerElem.current)
      }
    }
  }, [observerElem, handleObserver])

  if (isLoading && !data) {
    return (
      <div className="flex items-center justify-center h-screen">
        Loading...
      </div>
    )
  }
  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-red-500">Error loading data: {error.message}</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col md:max-w-none max-w-md mx-auto">
      <DesktopHeader />

      <div className="bg-white px-6 pb-4 md:px-8">
        <h1 className="text-xl font-bold text-gray-900 md:text-2xl md:mt-6">
          <span className="md:hidden">Welcome User!!</span>
          <span className="hidden md:block">Dashboard</span>
        </h1>
      </div>

      <div className="flex-1 px-6 pb-10 space-y-6 md:px-8 md:max-w-6xl md:mx-auto md:w-full">
        {/* Active session card - only show on mobile */}
        {activeSchedule && (
          <div className="md:hidden">
            <ActiveSessionCard
              customerName={activeSchedule.client_name}
              location={activeSchedule.location}
              time={activeSchedule.shiftTimeOnly}
              duration={duration}
              scheduleId={activeSchedule.id}
            />
          </div>
        )}

        {/* Stats cards - responsive layout */}
        <div className="space-y-4 md:space-y-6">
          <div className="md:hidden">
            <StatsCard
              value={stats.missed}
              label="Missed Scheduled"
              variant="error"
            />
          </div>

          <div className="flex gap-4 md:grid md:grid-cols-3 md:gap-8">
            <div className="hidden md:block">
              <StatsCard
                value={stats.missed}
                label="Missed Scheduled"
                variant="error"
                className="md:h-24 md:justify-center"
              />
            </div>
            <StatsCard
              value={stats.upcoming}
              label="Upcoming Today's Schedule"
              variant="warning"
              className="md:h-24 md:justify-center"
            />
            <StatsCard
              value={stats.completed}
              label="Today's Completed Schedule"
              variant="success"
              className="md:h-24 md:justify-center"
            />
          </div>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <h2 className="text-lg font-semibold text-gray-900">Schedule</h2>
              <div className="bg-cyan-500 text-white text-sm font-semibold px-2 py-0.5 rounded-lg">
                {allSchedules.length} {/* Display count of loaded schedules */}
              </div>
            </div>
            <button className="text-slate-700 text-sm hover:underline md:hidden">
              See All
            </button>
          </div>

          <div className="space-y-4">
            {allSchedules.length > 0 ? (
              allSchedules.map((sched) => (
                <AppointmentCard key={sched.id} schedule={sched} />
              ))
            ) : (
              <p className="text-center text-gray-500">No schedules found.</p>
            )}
            {/* Loading indicator for infinite scroll */}
            {isFetchingNextPage && (
              <div className="text-center text-gray-500 py-4">
                Loading more schedules...
              </div>
            )}
            {/* Observer element */}
            <div ref={observerElem} style={{ height: '10px' }} />
            {!hasNextPage && allSchedules.length > 0 && (
              <p className="text-center text-gray-500 py-4">
                You've reached the end of the list.
              </p>
            )}
          </div>
        </div>
      </div>

      {/* Footer for desktop */}
      <div className="hidden md:block bg-white py-4 text-center text-sm text-gray-600 border-t">
        ©2025 Careviah, Inc. All rights reserved.
      </div>

      {/* Bottom navigation - only show on mobile */}
      <div className="md:hidden">
        <BottomNavigation />
      </div>
    </div>
  )
}

function App() {
  return <Dashboard />
}
