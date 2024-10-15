package com.anxcye.utils;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.UUID;

/**
 * JWT工具类
 */
@Component
public class JwtUtil {

    private static Long JWT_TTL;
    private static String JWT_KEY;

    @Value ("${ancy.jwt.jwtExpiration}")
    public void setJwtExpiration(Long jwtExpiration) {
        JwtUtil.JWT_TTL = jwtExpiration * 1000L * 60 * 60 * 24;
    }

    @Value ("${ancy.jwt.jwtSecret}")
    public void setJwtKey(String jwtSecret) {
        JwtUtil.JWT_KEY = jwtSecret;
    }




    // 有效期为
//    public static Long JWT_TTL = 24 * 60 * 60 * 1000L * jwtExpiration;// 60 * 60 *1000 一个小时
    // 设置秘钥明文
//    public static final String JWT_KEY = "anxJWT signature does not match locally computed signature. JWT validity cannot be asserted and should not be trusted.cye";

    public static String getUUID() {
        return UUID.randomUUID().toString().replaceAll("-", "");
    }

    /**
     * 生成jtw
     *
     * @param subject token中要存放的数据（json格式）
     */
    public static String createJWT(String subject) {
        return getJwtBuilder(subject, null, getUUID());
    }

    /**
     * 生成jtw
     *
     * @param subject   token中要存放的数据（json格式）
     * @param ttlMillis token超时时间
     */
    public static String createJWT(String subject, Long ttlMillis) {
        return getJwtBuilder(subject, ttlMillis, getUUID());// 设置过期时间
    }

    private static String getJwtBuilder(String subject, Long ttlMillis, String uuid) {
        long nowMillis = System.currentTimeMillis();
        Date now = new Date(nowMillis);
        if (ttlMillis == null) {
            ttlMillis = JwtUtil.JWT_TTL;
        }
        long expMillis = nowMillis + ttlMillis;
        Date expDate = new Date(expMillis);
        return Jwts.builder()
                .subject(subject)
                // .issuer("ancy")
                .issuedAt(now)
                .signWith(generalKey())
                .expiration(expDate)
                .compact();
    }

    /**
     * 创建token
     */
    public static String createJWT(String id, String subject, Long ttlMillis) {
        return getJwtBuilder(subject, ttlMillis, id);// 设置过期时间
    }

    /**
     * 生成加密后的秘钥 secretKey
     */

    public static SecretKey generalKey() {
        byte[] encodedKey = JWT_KEY.getBytes(StandardCharsets.UTF_8);
        return new SecretKeySpec(encodedKey, 0, encodedKey.length, "HmacSHA256");
    }

    /**
     * 解析
     */
    public static Claims parseJWT(String jwt) throws Exception {
        return Jwts.parser()
                .verifyWith(generalKey())
                .build()
                .parseSignedClaims(jwt)
                .getPayload();
    }
}