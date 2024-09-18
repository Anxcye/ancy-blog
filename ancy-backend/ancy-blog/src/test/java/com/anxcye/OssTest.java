package com.anxcye;

import com.anxcye.properties.AliyunOssProperties;
import com.anxcye.utils.AliOssUtil;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;

@SpringBootTest(classes = AncyBlogApplication.class)
//@ExtendWith(MockitoSettings.class)
public class OssTest {

    @Autowired
    private AliOssUtil aliOssUtil;

    @Autowired
    private AliyunOssProperties aliyunOssProperties;

    @Test
    public void testOssProperties(){
        System.out.println(aliyunOssProperties.getAccessKeyId());
        System.out.println(aliyunOssProperties.getAccessKeySecret());
        System.out.println(aliyunOssProperties.getBucketName());
        System.out.println(aliyunOssProperties.getEndpoint());

    }


    @Test
    public void testOss() throws IOException {
//        AliyunOssProperties aliyunOssProperties = new AliyunOssProperties();
//        AliOssUtil aliOssUtil = new AliOssUtil(
//                aliyunOssProperties.getEndpoint(),
//                aliyunOssProperties.getAccessKeyId(),
//                aliyunOssProperties.getAccessKeySecret(),
//                aliyunOssProperties.getBucketName());
//
        //  /home/axy/Downloads/test.png
        String filePath = "/home/axy/Downloads/test.png";
        File file = new File(filePath);
        byte[] bytes = Files.readAllBytes(file.toPath());


        String result = aliOssUtil.upload(bytes, "test.jpg");
        System.out.println(result);

    }
}
