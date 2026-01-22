'use client'

interface StockBadgeProps {
  quantity: number
}

export function StockBadge({ quantity }: StockBadgeProps) {
  let bgColor = 'bg-green-100 text-green-800'
  let dotColor = 'bg-green-500'

  if (quantity <= 5) {
    bgColor = 'bg-red-100 text-red-800'
    dotColor = 'bg-red-500'
  } else if (quantity <= 10) {
    bgColor = 'bg-yellow-100 text-yellow-800'
    dotColor = 'bg-yellow-500'
  }

  return (
    <span
      className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-sm font-medium ${bgColor}`}
    >
      <span className={`w-2 h-2 rounded-full ${dotColor}`}></span>
      {quantity}
    </span>
  )
}
