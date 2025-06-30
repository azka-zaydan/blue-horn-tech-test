import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/profile')({
  component: Profile,
})

import { Button } from '@/components/ui/button'
import { BottomNavigation } from '@/components/dashboard/BottomNavigation'
import { Wifi, Signal, Battery } from 'lucide-react'

function StatusBar() {
  return (
    <div className="bg-white flex justify-between items-center px-6 py-4">
      <div className="text-gray-600 font-medium">9:41</div>
      <div className="flex items-center gap-1">
        <Signal className="h-5 w-5 text-gray-600" />
        <Wifi className="h-4 w-4 text-gray-600" />
        <Battery className="h-6 w-6 text-gray-600" />
      </div>
    </div>
  )
}

export default function Profile() {
  const handleLogOut = () => {
    console.log('Logging out...')
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col max-w-md mx-auto">
      <StatusBar />

      <div className="bg-white px-6 pb-4">
        <h1 className="text-xl font-bold text-gray-900">Welcome Louis!</h1>
      </div>

      <div className="flex-1 px-6 py-6">
        <Button
          onClick={handleLogOut}
          variant="outline"
          className="w-full border-red-300 text-red-600 hover:bg-red-50 py-4 text-base font-medium"
        >
          Log Out
        </Button>
      </div>

      <BottomNavigation />
    </div>
  )
}
