// src/components/dashboard/BottomNavigation.tsx
import React from 'react'
import { Home, User } from 'lucide-react' // Added Calendar and Settings for consistency
import { cn } from '@/lib/utils' // Assuming this utility is available
import { useNavigate, useLocation } from '@tanstack/react-router'

interface NavItemProps {
  icon: React.ComponentType<{ className?: string }>
  label: string
  active?: boolean
  onClick?: () => void
}

function NavItem({ icon: Icon, label, active = false, onClick }: NavItemProps) {
  return (
    <button
      onClick={onClick}
      className="flex flex-col items-center gap-1 flex-1 pb-1 hover:bg-blue-100 rounded-lg p-2 transition-colors"
    >
      <Icon
        className={cn('h-5 w-5', active ? 'text-slate-700' : 'text-gray-500')}
      />
      <span
        className={cn(
          'text-xs',
          active ? 'text-slate-700 font-semibold' : 'text-gray-500',
        )}
      >
        {label}
      </span>
    </button>
  )
}

export function BottomNavigation() {
  const navigate = useNavigate()
  const location = useLocation()

  const isHome = location.pathname === '/'
  const isProfile = location.pathname === '/profile'

  return (
    // Fixed position at the bottom, full width, z-index to stay on top
    // Adjusted styling to match the provided code's container (bg-blue-50, px-6 pt-4 pb-2)
    <nav className="fixed bottom-0 left-0 right-0 bg-blue-50 shadow-lg border-t border-gray-200 z-50 md:hidden">
      <div className="flex items-center gap-2 px-6 pt-4 pb-2 max-w-md mx-auto">
        {' '}
        {/* Added px, pt, pb from outer div */}
        <NavItem
          icon={Home}
          label="Home"
          active={isHome}
          onClick={() => navigate({ to: '/' })}
        />
        <NavItem
          icon={User}
          label="Profile"
          active={isProfile}
          onClick={() => navigate({ to: '/profile' })}
        />
      </div>
      <div className="flex justify-center mt-4 mb-2">
        {' '}
        {/* Added mb-2 for better spacing */}
        <div className="w-32 h-1 bg-gray-500 rounded-full" />
      </div>
    </nav>
  )
}
