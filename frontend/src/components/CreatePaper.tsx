import {Button, Input, message, Modal, Progress, Select, Space, Upload} from "antd";
import React, {useState} from "react";
import {LoadingOutlined, PlusOutlined} from "@ant-design/icons";
import {queryFile, queryFileURL, startUploadFile, uploadChunk} from "../apis/file";
import getFileHash from "../utils/hash";
import {createPapers} from "../apis/paper";

const SegmentSize = 1024 * 2

interface CreatePaperProps {
  opened: boolean
  onOk: () => void
  onCancel: () => void
}

function CreatePaper(props: CreatePaperProps) {

  const [
    loading,
    setLoading
  ] = React.useState<boolean>(false)

  const [
    uploaded, setUploaded
  ] = React.useState(false)

  const [
    percent, setPercent
  ] = React.useState(0)

  const [
    hash, setHash
  ] = React.useState('')

  const [
    emailTo, setEmailTo
  ] = React.useState('')

  const [
    language, setLanguage
  ] = React.useState('英语')

  const [
    messageApi,
    contextHolder
  ] = message.useMessage();

  return (
    <Modal open={props.opened}
           onOk={async () => {
             try {
               if (emailTo) {
                 if (!/^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/.test(emailTo)) {
                   messageApi.error("非法邮箱😁")
                   return
                 }
               }

               await createPapers(hash, language, emailTo)
               messageApi.info("创建翻译任务成功")
               setLoading(false)
               setUploaded(false)
               setPercent(0)
               setHash("")
               setEmailTo("")
               props.onOk()
             } catch (e) {
               messageApi.error("翻译失败")
             }
           }}
           onCancel={() => {
             setLoading(false)
             setUploaded(false)
             setPercent(0)
             setHash("")
             setEmailTo("")
             props.onCancel()
           }}
           okText={'开始翻译'}
           okButtonProps={{disabled: !uploaded}}
           cancelText={'取消'}
           title={'创建论文翻译'}
           destroyOnClose={true}
    >
      {contextHolder}
      {!uploaded ?
        <Upload.Dragger
          listType="text"
          showUploadList={false}
          customRequest={async (e) => {
            const file = e.file as File
            if (!(file.name.endsWith(".pdf") || file.name.endsWith(".PDF"))) {
              messageApi.error("只能上传pdf")
              return
            }
            if (file.size > 1024 * 1024 * 100) {
              messageApi.error("文件不能超过100M")
              return
            }
            setLoading(true)
            setUploaded(false)

            const hash = await getFileHash(file)
            setHash(hash)

            let startChunk = 0
            try {
              const fileInfo = await queryFile(hash)
              if (fileInfo.status !== 2) {
                startChunk = fileInfo.current_index
              } else {
                setLoading(false)
                setUploaded(true)
                setPercent(100)
                return
              }
            } catch (e) {
              startChunk = 0
            }

            const page = Math.ceil(file.size / SegmentSize);
            await startUploadFile(hash, file.name, page, SegmentSize)
            let start = 0;
            const promiseList = [];
            for (let i = startChunk; i < page; i++) {
              const end = start + SegmentSize;
              const item = file.slice(start, end);
              start = end;
              promiseList.push(uploadChunk(item, hash, i).then(() => {
                setPercent(parseInt((percent + (100 / page)).toFixed(0)))
              }))
            }

            try {
              await Promise.all(promiseList)
              setLoading(false)
              setUploaded(true)
              setPercent(100)
              const fileInfo = await queryFile(hash)
              if (fileInfo.status === 2) {
                const {url} = await queryFileURL(hash)
                console.log(url)
              }
            } catch (e) {
              messageApi.error("上传失败")
            }
          }}
        >
          {!loading ? (<div>
            <PlusOutlined/>
            上传PDF
          </div>) : (
            <div>
              <LoadingOutlined/>
              上传中
            </div>
          )}
        </Upload.Dragger> : ""
      }
      <div>
        {loading || uploaded ? <div>
          <Progress percent={percent}></Progress>
          <Space>
            <Select
              style={{width: 100}}
              defaultValue={language}
              onChange={(value) => {
                setLanguage(value)
              }}
              options={[
                {value: '英语', label: '英语'},
                {value: '汉语', label: '汉语'},
                {value: '韩语', label: '韩语'},
                {value: '日语', label: '日语'},
                {value: '法语', label: '法语'},
                {value: '德语', label: '德语'},
                {value: '西班牙语', label: '西班牙语'},
              ]}/>
            <Input placeholder={'通知邮箱'} onChange={(e) => {
              setEmailTo(e.target.value)
            }}/>
          </Space>

        </div> : ''}
      </div>
    </Modal>
  );
}

export default CreatePaper;
