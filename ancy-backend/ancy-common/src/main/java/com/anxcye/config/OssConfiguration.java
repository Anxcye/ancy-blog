package com.anxcye.config;

import com.anxcye.properties.AliyunOssProperties;
import com.anxcye.utils.AliOssUtil;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OssConfiguration {

    @Bean
    @ConditionalOnMissingBean
    public AliOssUtil aliOssUtil(AliyunOssProperties aliyunOssProperties){
        return new AliOssUtil(aliyunOssProperties.getEndpoint(),
                aliyunOssProperties.getAccessKeyId(),
                aliyunOssProperties.getAccessKeySecret(),
                aliyunOssProperties.getBucketName());

    }
}
