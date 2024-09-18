package com.anxcye.utils;

import com.aliyun.oss.OSS;
import com.aliyun.oss.OSSClientBuilder;
import com.aliyun.oss.OSSException;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.exception.SystemException;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;

import java.io.ByteArrayInputStream;

@Data
@AllArgsConstructor
@Slf4j
public class AliOssUtil {

    private String endpoint;
    private String accessKeyId;
    private String accessKeySecret;
    private String bucketName;

    /**
     * 文件上传
     *
     * @param bytes
     * @param objectName
     * @return
     */
    public String upload(byte[] bytes, String objectName) {

        // 创建OSSClient实例。
        OSS ossClient = new OSSClientBuilder().build(endpoint, accessKeyId, accessKeySecret);

        try {
            // 创建PutObject请求。
            ossClient.putObject(bucketName, objectName, new ByteArrayInputStream(bytes));
            // 文件访问路径规则 https://BucketName.Endpoint/ObjectName
            StringBuilder stringBuilder = new StringBuilder("https://");
            stringBuilder
                    .append(bucketName)
                    .append(".")
                    .append(endpoint)
                    .append("/")
                    .append(objectName);

            log.info("文件上传到:{}", stringBuilder);

            return stringBuilder.toString();
        } catch (OSSException oe) {
            throw new SystemException(AppHttpCodeEnum.UPLOAD_ERROR);
        } finally {
            if (ossClient != null) {
                ossClient.shutdown();
            }
        }

    }
}
