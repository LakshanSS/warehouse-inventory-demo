'use client'

import { AlertTriangle, AlertCircle } from 'lucide-react'

interface StockIndicatorProps {
  lowStock: number
  critical: number
}

export function StockIndicator({ lowStock, critical }: StockIndicatorProps) {
  return (
    <div className="flex flex-col sm:flex-row gap-4">
      {critical > 0 && (
        <div className="flex items-center gap-3 bg-red-50 border border-red-200 rounded-lg px-4 py-3">
          <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0" />
          <div>
            <p className="font-medium text-red-800">Critical Stock Alert</p>
            <p className="text-sm text-red-600">
              {critical} item{critical !== 1 ? 's' : ''} with 5 or fewer units
            </p>
          </div>
        </div>
      )}
      {lowStock > critical && (
        <div className="flex items-center gap-3 bg-yellow-50 border border-yellow-200 rounded-lg px-4 py-3">
          <AlertTriangle className="h-5 w-5 text-yellow-600 flex-shrink-0" />
          <div>
            <p className="font-medium text-yellow-800">Low Stock Warning</p>
            <p className="text-sm text-yellow-600">
              {lowStock - critical} item{lowStock - critical !== 1 ? 's' : ''}{' '}
              with 10 or fewer units
            </p>
          </div>
        </div>
      )}
    </div>
  )
}
