'use client'

import { useState } from 'react'
import { Product, CreateProductRequest } from '@/lib/types'

interface ProductFormProps {
  product?: Product
  onSubmit: (data: CreateProductRequest) => Promise<void>
  onCancel: () => void
  isLoading: boolean
}

export function ProductForm({ product, onSubmit, onCancel, isLoading }: ProductFormProps) {
  const isEditMode = !!product
  const [formData, setFormData] = useState<CreateProductRequest>({
    sku: product?.sku || '',
    name: product?.name || '',
    quantity: product?.quantity || 0,
    category: product?.category || '',
    location: product?.location || '',
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (!isEditMode) {
      if (!formData.sku.trim()) {
        newErrors.sku = 'SKU is required'
      }
      if (!formData.name.trim()) {
        newErrors.name = 'Name is required'
      }
    }

    if (formData.quantity !== undefined && formData.quantity < 0) {
      newErrors.quantity = 'Quantity must be 0 or greater'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    try {
      await onSubmit(formData)
    } catch {
      // Error is handled by parent
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* SKU */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          SKU {!isEditMode && <span className="text-red-500">*</span>}
        </label>
        {isEditMode ? (
          <div className="px-4 py-2 bg-gray-100 rounded-lg text-gray-700">
            {formData.sku}
          </div>
        ) : (
          <>
            <input
              type="text"
              value={formData.sku}
              onChange={(e) => setFormData({ ...formData, sku: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="e.g., WH-001"
              disabled={isLoading}
            />
            {errors.sku && <p className="text-sm text-red-600 mt-1">{errors.sku}</p>}
          </>
        )}
      </div>

      {/* Name */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Name {!isEditMode && <span className="text-red-500">*</span>}
        </label>
        {isEditMode ? (
          <div className="px-4 py-2 bg-gray-100 rounded-lg text-gray-700">
            {formData.name}
          </div>
        ) : (
          <>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Product name"
              disabled={isLoading}
            />
            {errors.name && <p className="text-sm text-red-600 mt-1">{errors.name}</p>}
          </>
        )}
      </div>

      {/* Quantity */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Quantity
        </label>
        <input
          type="number"
          value={formData.quantity}
          onChange={(e) => setFormData({ ...formData, quantity: parseInt(e.target.value) || 0 })}
          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          min="0"
          disabled={isLoading}
        />
        {errors.quantity && <p className="text-sm text-red-600 mt-1">{errors.quantity}</p>}
      </div>

      {/* Category - only in create mode */}
      {!isEditMode && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Category
          </label>
          <input
            type="text"
            value={formData.category}
            onChange={(e) => setFormData({ ...formData, category: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="e.g., Electronics"
            disabled={isLoading}
          />
        </div>
      )}

      {/* Location - only in create mode */}
      {!isEditMode && (
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Location
          </label>
          <input
            type="text"
            value={formData.location}
            onChange={(e) => setFormData({ ...formData, location: e.target.value })}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="e.g., Warehouse A, Shelf 1"
            disabled={isLoading}
          />
        </div>
      )}

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          disabled={isLoading}
          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isLoading}
          className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
        >
          {isLoading && (
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
          )}
          {isEditMode ? 'Save' : 'Create'}
        </button>
      </div>
    </form>
  )
}
