import React, { useState, useEffect } from 'react'; // Import useEffect
import { Card, Form, Input, Button, InputNumber, Typography, Select, message } from 'antd'; // Import message
import axios from 'axios'; // Import axios

const { Title } = Typography;
const { Option } = Select;

function TransactionSender() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [availableWallets, setAvailableWallets] = useState([]); // State for available wallets

  // TODO: Implement fetching available wallets on component mount
  useEffect(() => {
    // Placeholder for fetching wallets
    setAvailableWallets([
      { address: 'AZT123...', balance: 100 },
      { address: 'AZT456...', balance: 50 },
    ]);
  }, []);


  const onFinish = async (values) => { // Make function async
    setLoading(true);
    try {
      const transactionData = { // Prepare transaction data
        From: values.fromAddress,
        To: values.toAddress,
        Amount: values.amount,
        // Signature will be added later in backend or crypto module
      };
      const response = await axios.post('http://localhost:8080/transactions', transactionData); // Call create transaction API with data
      console.log('Transaction response:', response.data);
      message.success('交易发送成功 (占位符)'); // Use Ant Design message
      form.resetFields();
    } catch (error) {
      console.error('Error sending transaction:', error);
      message.error('交易发送失败'); // Use Ant Design message
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card>
      <Title level={2}>发起交易</Title>
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
      >
        <Form.Item
          label="发送方钱包"
          name="fromAddress"
          rules={[{ required: true, message: '请选择发送方钱包!' }]}
        >
          <Select placeholder="选择发送方钱包">
            {availableWallets.map(wallet => (
              <Option key={wallet.address} value={wallet.address}>
                {wallet.address} (余额: {wallet.balance})
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="接收方地址"
          name="toAddress"
          rules={[{ required: true, message: '请输入接收方地址!' }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="金额"
          name="amount"
          rules={[{ required: true, message: '请输入交易金额!' }, { type: 'number', min: 0, message: '金额必须大于0!' }]}
        >
          <InputNumber style={{ width: '100%' }} min={0} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading}>
            发送交易
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
}

export default TransactionSender;