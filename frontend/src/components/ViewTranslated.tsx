import {Button, message, Modal, Progress, Upload} from "antd";
import React, {useEffect, useState} from "react";
import {LoadingOutlined, PlusOutlined} from "@ant-design/icons";
import {queryFile, queryFileURL, startUploadFile, uploadChunk} from "../apis/file";
import getFileHash from "../utils/hash";
import {createPapers, getPaperDownloadURL, queryPaper} from "../apis/paper";
import TextArea from "antd/es/input/TextArea";

const SegmentSize = 1024 * 2

interface ViewTranslatedProps {
  opened: boolean
  paperID: string
  onOk: () => void
  onCancel: () => void
}

function ViewTranslated(props: ViewTranslatedProps) {


  const [
    messageApi,
    contextHolder
  ] = message.useMessage();

  const [
    docText,
    setDocText
  ] = React.useState('')

  const fetchPaper = async () => {
    const paper = await queryPaper(props.paperID)
    setDocText(paper.resultText)
  }

  useEffect(() => {
    fetchPaper()
  }, [props.paperID])

  return (
    <Modal open={props.opened}
           onOk={async () => {
             const paperDownloadURL = getPaperDownloadURL(props.paperID);
             const link = document.createElement('a');
             link.style.display = 'none';
             link.href = paperDownloadURL;
             link.setAttribute('download', "doc.pdf");
             document.body.appendChild(link);
             link.click();
             document.body.removeChild(link);
             props.onOk()
           }}
           onCancel={props.onCancel}
           okText={'下载'}
           cancelText={'取消'}
           title={'论文翻译后内容'}
           width={1000}
    >
      {contextHolder}
      <TextArea autoSize={true} value={docText}></TextArea>
    </Modal>
  );
}

export default ViewTranslated;
