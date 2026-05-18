import { View, Text } from '@tarojs/components'
import { useState, useEffect } from 'react'
import api from '../../services/api'
import './user.scss'

export default function Wallet() {
  const [wallet, setWallet] = useState<any>({})
  const [transactions, setTransactions] = useState<any[]>([])

  useEffect(() => { loadData() }, [])

  const loadData = async () => {
    try {
      const [wRes, tRes] = await Promise.all([api.walletBalance(), api.walletTransactions()])
      if (wRes.code === 0) setWallet(wRes.data)
      if (tRes.code === 0) setTransactions(tRes.data || [])
    } catch {}
  }

  const typeLabels: Record<string, string> = { deposit: '充值', withdraw: '提现', escrow_hold: '冻结', escrow_release: '放款', escrow_refund: '退款', commission: '佣金' }

  return (
    <View className='page-container'>
      <View className='wallet-card'>
        <Text className='wallet-label'>账户余额</Text>
        <Text className='wallet-balance'>¥{wallet.balance || '0.00'}</Text>
        <View className='wallet-sub'>
          <Text>冻结: ¥{wallet.frozen_amount || '0.00'}</Text>
        </View>
      </View>
      <Text className='section-title' style={{marginTop:'32px',marginBottom:'16px',display:'block'}}>交易记录</Text>
      {transactions.map(t => (
        <View key={t.id} className='menu-item'>
          <View style={{flex:1}}>
            <Text style={{fontSize:'28px',color:'#303133'}}>{typeLabels[t.type] || t.type}</Text>
            <Text style={{fontSize:'22px',color:'#909399',display:'block',marginTop:'4px'}}>{t.description || ''}</Text>
          </View>
          <Text style={{fontSize:'30px',color:t.amount>0?'#67C23A':'#F56C6C',fontWeight:'bold'}}>
            {t.amount > 0 ? '+' : ''}¥{t.amount}
          </Text>
        </View>
      ))}
    </View>
  )
}
