'use client'

import { Pencil, Trash2 } from 'lucide-react'
import { Product } from '@/lib/types'

interface ActionButtonsProps {
  product: Product
  onEdit: (product: Product) => void
  onDelete: (product: Product) => void
}

export function ActionButtons({ product, onEdit, onDelete }: ActionButtonsProps) {
  return (
    <div className="flex items-center gap-2">
      <button
        onClick={() => onEdit(product)}
        className="p-1.5 rounded hover:bg-gray-100 transition-colors"
        title="Edit quantity"
      >
        <Pencil className="h-4 w-4 text-gray-500 hover:text-blue-600" />
      </button>
      <button
        onClick={() => onDelete(product)}
        className="p-1.5 rounded hover:bg-gray-100 transition-colors"
        title="Delete product"
      >
        <Trash2 className="h-4 w-4 text-gray-500 hover:text-red-600" />
      </button>
    </div>
  )
}
