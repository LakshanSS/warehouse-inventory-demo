'use client'

import { useEffect, useState } from 'react'
import { Package, AlertTriangle, RefreshCw, Search } from 'lucide-react'
import { InventoryTable } from '@/components/InventoryTable'
import { StockIndicator } from '@/components/StockIndicator'
import { Product } from '@/lib/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090'

export default function Home() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')

  const fetchInventory = async () => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`${API_URL}/inventory`)
      if (!response.ok) throw new Error('Failed to fetch inventory')
      const data = await response.json()
      setProducts(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchInventory()
  }, [])

  const lowStockCount = products.filter(p => p.quantity <= 10).length
  const criticalStockCount = products.filter(p => p.quantity <= 5).length
  const totalItems = products.reduce((sum, p) => sum + p.quantity, 0)

  const filteredProducts = products.filter(
    p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.sku.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.category?.toLowerCase().includes(searchQuery.toLowerCase())
  )

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
            <Package className="h-8 w-8 text-blue-600" />
            Warehouse Inventory Dashboard
          </h1>
          <p className="text-gray-600 mt-2">
            Monitor and manage your warehouse stock levels in real-time
          </p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-sm font-medium text-gray-500">Total Products</div>
            <div className="text-2xl font-bold text-gray-900">{products.length}</div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-sm font-medium text-gray-500">Total Items in Stock</div>
            <div className="text-2xl font-bold text-gray-900">{totalItems.toLocaleString()}</div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-sm font-medium text-gray-500">Low Stock Items</div>
            <div className="text-2xl font-bold text-yellow-600 flex items-center gap-2">
              {lowStockCount}
              {lowStockCount > 0 && <AlertTriangle className="h-5 w-5" />}
            </div>
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <div className="text-sm font-medium text-gray-500">Critical Stock</div>
            <div className="text-2xl font-bold text-red-600 flex items-center gap-2">
              {criticalStockCount}
              {criticalStockCount > 0 && <AlertTriangle className="h-5 w-5" />}
            </div>
          </div>
        </div>

        {/* Stock Indicators */}
        {(lowStockCount > 0 || criticalStockCount > 0) && (
          <div className="mb-6">
            <StockIndicator
              lowStock={lowStockCount}
              critical={criticalStockCount}
            />
          </div>
        )}

        {/* Search and Refresh */}
        <div className="flex flex-col sm:flex-row gap-4 mb-6">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search by name, SKU, or category..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <button
            onClick={fetchInventory}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            <RefreshCw className={`h-5 w-5 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </button>
        </div>

        {/* Error State */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        {/* Table */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <InventoryTable
            data={filteredProducts}
            loading={loading}
            onRefresh={fetchInventory}
          />
        </div>
      </div>
    </main>
  )
}
